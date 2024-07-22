package component

import (
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/unchainio/pkg/errors"
)

var NotAllowedErr = sql.ErrNoRows

func must(flag bool) error {
	if !flag {
		return errors.Wrap(NotAllowedErr, "")
	}
	return nil
}

func canUseBase(org *orm.Organization, base *orm.Base) bool {
	// TODO if it's a paid public component, without an active subscription, make it unavailable
	return canViewBase(org, base)
}

func canUseTrigger(org *orm.Organization, trigger *orm.Trigger) bool {
	// TODO if it's a paid public component, without an active subscription, make it unavailable
	return canViewTrigger(org, trigger)
}

func canUseAction(org *orm.Organization, action *orm.Action) bool {
	// TODO if it's a paid public component, without an active subscription, make it unavailable
	return canViewAction(org, action)
}

func canViewBase(org *orm.Organization, base *orm.Base) bool {
	return base.Public || canEditBase(org, base)
}

func canViewTrigger(org *orm.Organization, trigger *orm.Trigger) bool {
	return trigger.Public || canEditTrigger(org, trigger)
}

func canViewAction(org *orm.Organization, action *orm.Action) bool {
	return action.Public || canEditAction(org, action)
}

func canEditBase(org *orm.Organization, base *orm.Base) bool {
	return base.DeveloperID == org.ID
}

func canEditTrigger(org *orm.Organization, trigger *orm.Trigger) bool {
	return trigger.DeveloperID == org.ID
}

func canEditAction(org *orm.Organization, action *orm.Action) bool {
	return action.DeveloperID == org.ID
}

func canViewBaseVersion(org *orm.Organization, base *orm.Base, baseVersion *orm.BaseVersion) bool {
	return (base.Public && baseVersion.Public) || canEditBase(org, base)
}

func canViewTriggerVersion(org *orm.Organization, trigger *orm.Trigger, triggerVersion *orm.TriggerVersion) bool {
	return (trigger.Public && triggerVersion.Public) || canEditTrigger(org, trigger)
}

func canViewActionVersion(org *orm.Organization, action *orm.Action, actionVersion *orm.ActionVersion) bool {
	return (action.Public && actionVersion.Public) || canEditAction(org, action)
}

func canUseBaseVersion(org *orm.Organization, base *orm.Base, baseVersion *orm.BaseVersion) bool {
	// TODO if it's a paid public component, without an active subscription, make it unavailable
	return canViewBaseVersion(org, base, baseVersion)
}

func canUseTriggerVersion(org *orm.Organization, trigger *orm.Trigger, triggerVersion *orm.TriggerVersion) bool {
	// TODO if it's a paid public component, without an active subscription, make it unavailable
	return canViewTriggerVersion(org, trigger, triggerVersion)
}

func canUseActionVersion(org *orm.Organization, action *orm.Action, actionVersion *orm.ActionVersion) bool {
	// TODO if it's a paid public component, without an active subscription, make it unavailable
	return canViewActionVersion(org, action, actionVersion)
}
