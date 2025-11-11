/*
 * Jolt Physics C Wrapper - Body Operations Implementation
 */

#include "body.h"
#include "physics.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Body/Body.h>
#include <Jolt/Physics/Body/BodyCreationSettings.h>
#include <Jolt/Physics/Body/BodyInterface.h>
#include <memory>

using namespace JPH;

// Collision layers (defined in physics.cpp)
namespace Layers
{
	static constexpr ObjectLayer NON_MOVING = 0;
	static constexpr ObjectLayer MOVING = 1;
};

JoltBodyInterface JoltPhysicsSystemGetBodyInterface(JoltPhysicsSystem system)
{
	PhysicsSystemWrapper *wrapper = static_cast<PhysicsSystemWrapper *>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);
	BodyInterface* bi = &ps->GetBodyInterface();

	return static_cast<JoltBodyInterface>(bi);
}

void JoltGetBodyPosition(const JoltBodyInterface bodyInterface,
						 const JoltBodyID bodyID,
						 float *x, float *y, float *z)
{
	const BodyInterface *bi = static_cast<const BodyInterface *>(bodyInterface);
	const BodyID *bid = static_cast<const BodyID *>(bodyID);

	RVec3 pos = bi->GetPosition(*bid);
	*x = static_cast<float>(pos.GetX());
	*y = static_cast<float>(pos.GetY());
	*z = static_cast<float>(pos.GetZ());
}

JoltBodyID JoltCreateBody(JoltBodyInterface bodyInterface,
						  JoltShape shape,
						  float x, float y, float z,
						  int isDynamic)
{
	BodyInterface *bi = static_cast<BodyInterface *>(bodyInterface);
	const Shape *s = static_cast<const Shape *>(shape);

	BodyCreationSettings body_settings(
		s,
		RVec3(x, y, z),
		Quat::sIdentity(),
		isDynamic ? EMotionType::Dynamic : EMotionType::Static,
		isDynamic ? Layers::MOVING : Layers::NON_MOVING);

	Body *body = bi->CreateBody(body_settings);
	bi->AddBody(body->GetID(), EActivation::Activate);

	// Use smart pointer for exception safety, then release to caller
	auto bodyIDPtr = std::make_unique<BodyID>(body->GetID());
	return static_cast<JoltBodyID>(bodyIDPtr.release());
}

void JoltDestroyBodyID(JoltBodyID bodyID)
{
	BodyID *bid = static_cast<BodyID *>(bodyID);
	delete bid;
}
