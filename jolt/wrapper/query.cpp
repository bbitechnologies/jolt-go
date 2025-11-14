/*
 * Jolt Physics C Wrapper - Collision Query Implementation
 */

#include "query.h"
#include "physics.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Collision/CollideShape.h>
#include <Jolt/Physics/Collision/CollisionCollectorImpl.h>
#include <Jolt/Physics/Collision/NarrowPhaseQuery.h>
#include <Jolt/Physics/Body/BodyID.h>
#include <vector>

using namespace JPH;

// Collision layers (defined in physics.cpp)
namespace Layers
{
	static constexpr ObjectLayer MOVING = 1;
};

// Adapter: converts ObjectVsBroadPhaseLayerFilter to BroadPhaseLayerFilter for queries
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

// Adapter: converts ObjectLayerPairFilter to ObjectLayerFilter for queries
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

// Collector that just checks if any collision occurred
class AnyHitCollector : public CollideShapeCollector
{
public:
	AnyHitCollector() : m_hasHit(false) {}

	virtual void AddHit(const CollideShapeResult& inResult) override
	{
		m_hasHit = true;
	}

	bool HasHit() const { return m_hasHit; }

private:
	bool m_hasHit;
};

// Collector that stores all collision hits
class AllHitsCollector : public CollideShapeCollector
{
public:
	AllHitsCollector(JoltCollisionHit* outHits, int maxHits)
		: m_outHits(outHits), m_maxHits(maxHits), m_numHits(0) {}

	virtual void AddHit(const CollideShapeResult& inResult) override
	{
		if (m_numHits < m_maxHits)
		{
			JoltCollisionHit& hit = m_outHits[m_numHits];

			// Store body ID
			BodyID* bodyIDCopy = new BodyID(inResult.mBodyID2);
			hit.bodyID = static_cast<JoltBodyID>(bodyIDCopy);

			// Store contact point (using contact point on second shape)
			Vec3 contactPoint = inResult.mContactPointOn2;
			hit.contactPointX = contactPoint.GetX();
			hit.contactPointY = contactPoint.GetY();
			hit.contactPointZ = contactPoint.GetZ();

			// Store penetration depth
			hit.penetrationDepth = inResult.mPenetrationDepth;

			m_numHits++;
		}
	}

	int GetNumHits() const { return m_numHits; }

private:
	JoltCollisionHit* m_outHits;
	int m_maxHits;
	int m_numHits;
};

int JoltCollideShape(JoltPhysicsSystem system, JoltShape shape,
                     float posX, float posY, float posZ, float penetrationTolerance)
{
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);
	const Shape* s = static_cast<const Shape*>(shape);

	// Get narrow phase query interface
	const NarrowPhaseQuery& query = ps->GetNarrowPhaseQuery();

	// Create filter adapters (query shape acts as MOVING layer)
	BroadPhaseLayerFilterAdapter bpFilter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter objFilter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Create collector to check for any hit
	AnyHitCollector collector;

	// Setup collision settings
	CollideShapeSettings settings;
	settings.mPenetrationTolerance = penetrationTolerance;

	// Perform collision query
	query.CollideShape(
		s,
		Vec3::sReplicate(1.0f),  // Scale
		RMat44::sTranslation(RVec3(posX, posY, posZ)),  // Transform (position, no rotation)
		settings,
		RVec3::sZero(),  // Base offset
		collector,
		bpFilter,
		objFilter
	);

	return collector.HasHit() ? 1 : 0;
}

int JoltCollideShapeGetHits(JoltPhysicsSystem system, JoltShape shape,
                            float posX, float posY, float posZ,
                            JoltCollisionHit* outHits, int maxHits, float penetrationTolerance)
{
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);
	const Shape* s = static_cast<const Shape*>(shape);

	// Get narrow phase query interface
	const NarrowPhaseQuery& query = ps->GetNarrowPhaseQuery();

	// Create filter adapters (query shape acts as MOVING layer)
	BroadPhaseLayerFilterAdapter bpFilter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter objFilter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Create collector to gather all hits
	AllHitsCollector collector(outHits, maxHits);

	// Setup collision settings
	CollideShapeSettings settings;
	settings.mPenetrationTolerance = penetrationTolerance;

	// Perform collision query
	query.CollideShape(
		s,
		Vec3::sReplicate(1.0f),  // Scale
		RMat44::sTranslation(RVec3(posX, posY, posZ)),  // Transform (position, no rotation)
		settings,
		RVec3::sZero(),  // Base offset
		collector,
		bpFilter,
		objFilter
	);

	return collector.GetNumHits();
}
