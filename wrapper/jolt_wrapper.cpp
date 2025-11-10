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
#include <Jolt/Physics/Body/BodyCreationSettings.h>
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
