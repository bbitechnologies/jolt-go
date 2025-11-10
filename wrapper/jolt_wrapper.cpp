/*
 * Jolt Physics C Wrapper Implementation
 *
 * CRITICAL BUG FIX: PhysicsSystem::Init() takes references to layer interfaces,
 * so they must be heap-allocated and kept alive. We use PhysicsSystemWrapper
 * to manage their lifetime alongside the PhysicsSystem.
 */

#include "jolt_wrapper.h"
#include <Jolt/Jolt.h>
#include <Jolt/RegisterTypes.h>
#include <Jolt/Core/Factory.h>
#include <Jolt/Core/TempAllocator.h>
#include <Jolt/Core/JobSystemThreadPool.h>
#include <Jolt/Physics/PhysicsSettings.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Collision/Shape/SphereShape.h>
#include <Jolt/Physics/Collision/Shape/BoxShape.h>
#include <Jolt/Physics/Collision/Shape/CapsuleShape.h>
#include <Jolt/Physics/Body/BodyCreationSettings.h>
#include <Jolt/Physics/Character/CharacterVirtual.h>
#include <iostream>
#include <cstdarg>

using namespace JPH;

// Trace callback for Jolt error messages (helps with debugging)
static void TraceImpl(const char *inFMT, ...)
{
	va_list list;
	va_start(list, inFMT);
	char buffer[1024];
	vsnprintf(buffer, sizeof(buffer), inFMT, list);
	va_end(list);
	std::cout << buffer << std::endl;
}

#ifdef JPH_ENABLE_ASSERTS
static bool AssertFailedImpl(const char *inExpression, const char *inMessage, const char *inFile, uint inLine)
{
	std::cout << inFile << ":" << inLine << ": (" << inExpression << ") "
			  << (inMessage != nullptr ? inMessage : "") << std::endl;
	return true;
};
#endif

// Global Jolt resources (shared by all PhysicsSystems)
static TempAllocatorImpl *gTempAllocator = nullptr;
static JobSystemThreadPool *gJobSystem = nullptr;

int JoltInit()
{
	Trace = TraceImpl;
	JPH_IF_ENABLE_ASSERTS(AssertFailed = AssertFailedImpl;)

	Factory::sInstance = new Factory();
	RegisterTypes();

	gTempAllocator = new TempAllocatorImpl(10 * 1024 * 1024);
	gJobSystem = new JobSystemThreadPool(cMaxPhysicsJobs, cMaxPhysicsBarriers,
										 std::thread::hardware_concurrency() - 1);
	return 1;
}

void JoltShutdown()
{
	delete gJobSystem;
	delete gTempAllocator;
	delete Factory::sInstance;
	Factory::sInstance = nullptr;
}

// Collision layers: NON_MOVING (static) and MOVING (dynamic)
namespace Layers
{
	static constexpr ObjectLayer NON_MOVING = 0;
	static constexpr ObjectLayer MOVING = 1;

	static constexpr ObjectLayer NUM_LAYERS = 2;
};

namespace BroadPhaseLayers
{
	static constexpr BroadPhaseLayer NON_MOVING(0);
	static constexpr BroadPhaseLayer MOVING(1);

	static constexpr uint NUM_LAYERS(2);
};

// Maps object layers to broad phase layers
class BPLayerInterfaceImpl final : public BroadPhaseLayerInterface
{
public:
	BPLayerInterfaceImpl()
	{
		mObjectToBroadPhase[Layers::NON_MOVING] = BroadPhaseLayers::NON_MOVING;
		mObjectToBroadPhase[Layers::MOVING] = BroadPhaseLayers::MOVING;
	}

	virtual uint GetNumBroadPhaseLayers() const override
	{
		return BroadPhaseLayers::NUM_LAYERS;
	}

	virtual BroadPhaseLayer GetBroadPhaseLayer(ObjectLayer inLayer) const override
	{
		JPH_ASSERT(inLayer < Layers::NUM_LAYERS);
		return mObjectToBroadPhase[inLayer];
	}

#if defined(JPH_EXTERNAL_PROFILE) || defined(JPH_PROFILE_ENABLED)
	virtual const char* GetBroadPhaseLayerName(BroadPhaseLayer inLayer) const override
	{
		switch ((BroadPhaseLayer::Type)inLayer)
		{
		case (BroadPhaseLayer::Type)BroadPhaseLayers::NON_MOVING:	return "NON_MOVING";
		case (BroadPhaseLayer::Type)BroadPhaseLayers::MOVING:		return "MOVING";
		default:													return "INVALID";
		}
	}
#endif // JPH_EXTERNAL_PROFILE || JPH_PROFILE_ENABLED

private:
	BroadPhaseLayer mObjectToBroadPhase[Layers::NUM_LAYERS];
};

// Filters which broad phase layers can collide
class ObjectVsBroadPhaseLayerFilterImpl : public ObjectVsBroadPhaseLayerFilter
{
public:
	virtual bool ShouldCollide(ObjectLayer inLayer1, BroadPhaseLayer inLayer2) const override
	{
		switch (inLayer1)
		{
		case Layers::NON_MOVING:
			return inLayer2 == BroadPhaseLayers::MOVING;
		case Layers::MOVING:
			return true;
		default:
			JPH_ASSERT(false);
			return false;
		}
	}
};

// Filters which object layers can collide with each other
class ObjectLayerPairFilterImpl : public ObjectLayerPairFilter
{
public:
	virtual bool ShouldCollide(ObjectLayer inObject1, ObjectLayer inObject2) const override
	{
		switch (inObject1)
		{
		case Layers::NON_MOVING:
			return inObject2 == Layers::MOVING;
		case Layers::MOVING:
			return true;
		default:
			JPH_ASSERT(false);
			return false;
		}
	}
};

// Adapter: converts ObjectVsBroadPhaseLayerFilter to BroadPhaseLayerFilter for character collision
class BroadPhaseLayerFilterAdapter : public BroadPhaseLayerFilter
{
public:
	BroadPhaseLayerFilterAdapter(const ObjectVsBroadPhaseLayerFilterImpl* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(BroadPhaseLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectVsBroadPhaseLayerFilterImpl* m_filter;
	ObjectLayer m_object_layer;
};

// Adapter: converts ObjectLayerPairFilter to ObjectLayerFilter for character collision
class ObjectLayerFilterAdapter : public ObjectLayerFilter
{
public:
	ObjectLayerFilterAdapter(const ObjectLayerPairFilterImpl* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(ObjectLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectLayerPairFilterImpl* m_filter;
	ObjectLayer m_object_layer;
};

// Wrapper to keep layer interfaces alive (PhysicsSystem stores references to them)
struct PhysicsSystemWrapper
{
	PhysicsSystem *system;
	BPLayerInterfaceImpl *broad_phase_layer_interface;
	ObjectVsBroadPhaseLayerFilterImpl *object_vs_broadphase_layer_filter;
	ObjectLayerPairFilterImpl *object_vs_object_layer_filter;
};

JoltPhysicsSystem JoltCreatePhysicsSystem()
{
	const uint cMaxBodies = 1024;
	const uint cNumBodyMutexes = 0;
	const uint cMaxBodyPairs = 1024;
	const uint cMaxContactConstraints = 1024;

	// Heap-allocate layer interfaces (must outlive PhysicsSystem!)
	BPLayerInterfaceImpl *broad_phase_layer_interface = new BPLayerInterfaceImpl();
	ObjectVsBroadPhaseLayerFilterImpl *object_vs_broadphase_layer_filter = new ObjectVsBroadPhaseLayerFilterImpl();
	ObjectLayerPairFilterImpl *object_vs_object_layer_filter = new ObjectLayerPairFilterImpl();

	PhysicsSystem *system = new PhysicsSystem();
	system->Init(cMaxBodies, cNumBodyMutexes, cMaxBodyPairs, cMaxContactConstraints,
				 *broad_phase_layer_interface,
				 *object_vs_broadphase_layer_filter,
				 *object_vs_object_layer_filter);

	// Wrap everything for proper cleanup
	PhysicsSystemWrapper *wrapper = new PhysicsSystemWrapper();
	wrapper->system = system;
	wrapper->broad_phase_layer_interface = broad_phase_layer_interface;
	wrapper->object_vs_broadphase_layer_filter = object_vs_broadphase_layer_filter;
	wrapper->object_vs_object_layer_filter = object_vs_object_layer_filter;

	return static_cast<JoltPhysicsSystem>(wrapper);
}

void JoltDestroyPhysicsSystem(JoltPhysicsSystem system)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	delete wrapper->system;
	delete wrapper->broad_phase_layer_interface;
	delete wrapper->object_vs_broadphase_layer_filter;
	delete wrapper->object_vs_object_layer_filter;
	delete wrapper;
}

void JoltPhysicsSystemUpdate(JoltPhysicsSystem system, float deltaTime)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	wrapper->system->Update(deltaTime, 1, gTempAllocator, gJobSystem);
}

JoltBodyInterface JoltPhysicsSystemGetBodyInterface(JoltPhysicsSystem system)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	return static_cast<JoltBodyInterface>(&wrapper->system->GetBodyInterface());
}

void JoltGetBodyPosition(JoltBodyInterface bodyInterface,
						 JoltBodyID bodyID,
						 float *x, float *y, float *z)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	BodyID *bid = static_cast<BodyID *>(bodyID);

	RVec3 pos = bi->GetPosition(*bid);
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

JoltBodyID JoltCreateSphere(JoltBodyInterface bodyInterface,
							float radius,
							float x, float y, float z,
							int isDynamic)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);

	SphereShapeSettings sphere_settings(radius);
	ShapeSettings::ShapeResult sphere_result = sphere_settings.Create();

	BodyCreationSettings body_settings(
		sphere_result.Get(),
		RVec3(x, y, z),
		Quat::sIdentity(),
		isDynamic ? EMotionType::Dynamic : EMotionType::Static,
		isDynamic ? Layers::MOVING : Layers::NON_MOVING);

	Body *body = bi->CreateBody(body_settings);
	bi->AddBody(body->GetID(), EActivation::Activate);

	BodyID *bodyIDPtr = new BodyID(body->GetID());
	return static_cast<JoltBodyID>(bodyIDPtr);
}

JoltBodyID JoltCreateBox(JoltBodyInterface bodyInterface,
						 float halfExtentX, float halfExtentY, float halfExtentZ,
						 float x, float y, float z,
						 int isDynamic)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);

	BoxShapeSettings box_settings(Vec3(halfExtentX, halfExtentY, halfExtentZ));
	ShapeSettings::ShapeResult box_result = box_settings.Create();

	BodyCreationSettings body_settings(
		box_result.Get(),
		RVec3(x, y, z),
		Quat::sIdentity(),
		isDynamic ? EMotionType::Dynamic : EMotionType::Static,
		isDynamic ? Layers::MOVING : Layers::NON_MOVING);

	Body *body = bi->CreateBody(body_settings);
	bi->AddBody(body->GetID(), EActivation::Activate);

	BodyID *bodyIDPtr = new BodyID(body->GetID());
	return static_cast<JoltBodyID>(bodyIDPtr);
}

void JoltDestroyBodyID(JoltBodyID bodyID)
{
	BodyID *bid = static_cast<BodyID *>(bodyID);
	delete bid;
}

JoltCharacterVirtual JoltCreateCharacterVirtual(JoltPhysicsSystem system,
											 float x, float y, float z)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);

	// Create a capsule shape for the character (half-height 0.9m, radius 0.5m = ~1.8m tall human)
	CapsuleShapeSettings capsule_settings(0.9f, 0.5f);
	ShapeSettings::ShapeResult capsule_result = capsule_settings.Create();

	CharacterVirtualSettings settings;
	settings.mShape = capsule_result.Get();
	settings.mMaxSlopeAngle = DegreesToRadians(45.0f);
	settings.mMass = 70.0f; // 70kg
	settings.mMaxStrength = 100.0f;
	settings.mCharacterPadding = 0.02f;
	settings.mPenetrationRecoverySpeed = 1.0f;
	settings.mPredictiveContactDistance = 0.1f;

	// Create at specified position
	CharacterVirtual* character = new CharacterVirtual(&settings, RVec3(x, y, z), Quat::sIdentity(), wrapper->system);
	return static_cast<JoltCharacterVirtual>(character);
}

void JoltDestroyCharacterVirtual(JoltCharacterVirtual character)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	delete cv;
}

void JoltCharacterVirtualExtendedUpdate(JoltCharacterVirtual character,
										JoltPhysicsSystem system,
										float deltaTime,
										float gravityX, float gravityY, float gravityZ)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);

	// Use default extended update settings
	CharacterVirtual::ExtendedUpdateSettings settings;

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(wrapper->object_vs_broadphase_layer_filter, Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(wrapper->object_vs_object_layer_filter, Layers::MOVING);

	// Call ExtendedUpdate with gravity vector and layer filters
	cv->ExtendedUpdate(
		deltaTime,
		Vec3(gravityX, gravityY, gravityZ),
		settings,
		broad_phase_filter,
		object_layer_filter,
		{}, // Empty BodyFilter (collides with all bodies)
		{}, // Empty ShapeFilter (collides with all shapes)
		*gTempAllocator
	);
}

void JoltCharacterVirtualSetLinearVelocity(JoltCharacterVirtual character,
										   float x, float y, float z)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	cv->SetLinearVelocity(Vec3(x, y, z));
}

void JoltCharacterVirtualGetPosition(JoltCharacterVirtual character,
									 float* x, float* y, float* z)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	RVec3 pos = cv->GetPosition();
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

JoltGroundState JoltCharacterVirtualGetGroundState(JoltCharacterVirtual character)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	CharacterBase::EGroundState state = cv->GetGroundState();
	return static_cast<JoltGroundState>(state);
}

int JoltCharacterVirtualIsSupported(JoltCharacterVirtual character)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	return cv->IsSupported() ? 1 : 0;
}
