/*
 * Jolt Physics C Wrapper - Character Virtual Implementation
 */

#include "character.h"
#include "physics.h"
#include "core.h"
#include <Jolt/Jolt.h>
#include <Jolt/Core/TempAllocator.h>
#include <Jolt/Physics/Collision/Shape/CapsuleShape.h>
#include <Jolt/Physics/Character/CharacterVirtual.h>
#include <memory>

using namespace JPH;

// Collision layers (defined in physics.cpp)
namespace Layers
{
	static constexpr ObjectLayer MOVING = 1;
};

// Adapter: converts ObjectVsBroadPhaseLayerFilter to BroadPhaseLayerFilter for character collision
class BroadPhaseLayerFilterAdapter : public BroadPhaseLayerFilter
{
public:
	BroadPhaseLayerFilterAdapter(const ObjectVsBroadPhaseLayerFilter* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(BroadPhaseLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectVsBroadPhaseLayerFilter* m_filter;
	ObjectLayer m_object_layer;
};

// Adapter: converts ObjectLayerPairFilter to ObjectLayerFilter for character collision
class ObjectLayerFilterAdapter : public ObjectLayerFilter
{
public:
	ObjectLayerFilterAdapter(const ObjectLayerPairFilter* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(ObjectLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectLayerPairFilter* m_filter;
	ObjectLayer m_object_layer;
};

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

	// Create at specified position using smart pointer for exception safety
	auto character = std::make_unique<CharacterVirtual>(&settings, RVec3(x, y, z), Quat::sIdentity(), GetPhysicsSystem(wrapper));
	return static_cast<JoltCharacterVirtual>(character.release());
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
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Call ExtendedUpdate with gravity vector and layer filters
	cv->ExtendedUpdate(
		deltaTime,
		Vec3(gravityX, gravityY, gravityZ),
		settings,
		broad_phase_filter,
		object_layer_filter,
		{}, // Empty BodyFilter (collides with all bodies)
		{}, // Empty ShapeFilter (collides with all shapes)
		*gTempAllocator.get()
	);
}

void JoltCharacterVirtualSetLinearVelocity(JoltCharacterVirtual character,
										   float x, float y, float z)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	cv->SetLinearVelocity(Vec3(x, y, z));
}

void JoltCharacterVirtualGetLinearVelocity(const JoltCharacterVirtual character,
										   float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	Vec3 vel = cv->GetLinearVelocity();
	*x = vel.GetX();
	*y = vel.GetY();
	*z = vel.GetZ();
}

void JoltCharacterVirtualGetGroundVelocity(const JoltCharacterVirtual character,
											float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	Vec3 vel = cv->GetGroundVelocity();
	*x = vel.GetX();
	*y = vel.GetY();
	*z = vel.GetZ();
}

void JoltCharacterVirtualGetPosition(const JoltCharacterVirtual character,
									 float* x, float* y, float* z)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	RVec3 pos = cv->GetPosition();
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

JoltGroundState JoltCharacterVirtualGetGroundState(const JoltCharacterVirtual character)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	CharacterBase::EGroundState state = cv->GetGroundState();
	return static_cast<JoltGroundState>(state);
}

int JoltCharacterVirtualIsSupported(const JoltCharacterVirtual character)
{
	const CharacterVirtual* cv = static_cast<const CharacterVirtual*>(character);
	return cv->IsSupported() ? 1 : 0;
}

void JoltCharacterVirtualSetShape(JoltCharacterVirtual character,
								  JoltShape shape,
								  float maxPenetrationDepth,
								  JoltPhysicsSystem system)
{
	CharacterVirtual* cv = static_cast<CharacterVirtual*>(character);
	const Shape* s = static_cast<const Shape*>(shape);
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);

	// Use MOVING layer for character (same as dynamic bodies)
	BroadPhaseLayerFilterAdapter broad_phase_filter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter object_layer_filter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Call SetShape with required filters
	cv->SetShape(
		s,
		maxPenetrationDepth,
		broad_phase_filter,
		object_layer_filter,
		{}, // Empty BodyFilter (collides with all bodies)
		{}, // Empty ShapeFilter (collides with all shapes)
		*gTempAllocator.get()
	);
}
