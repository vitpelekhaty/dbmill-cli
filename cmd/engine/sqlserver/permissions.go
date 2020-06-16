package sqlserver

import (
	"fmt"
)

// Permissions разрешения
type Permissions map[string]bool

// List возвращает срез разрешений
func (perms Permissions) List() []string {
	if len(perms) == 0 {
		return nil
	}

	list := make([]string, len(perms))

	var index int

	for perm := range perms {
		list[index] = perm
		index++
	}

	return list
}

// PermissionState состояние разрешения
type PermissionState byte

const (
	// PermStateUnknown неизвестное состояние
	PermStateUnknown PermissionState = iota
	// PermStateGrant GRANT
	PermStateGrant
	// PermStateGrantWithGrantOption GRANT_WITH_GRANT_OPTION
	PermStateGrantWithGrantOption
	// PermStateDeny DENY
	PermStateDeny
	// PermStateRevoke REVOKE
	PermStateRevoke
)

// String строковое представление PermissionState
func (ps PermissionState) String() string {
	switch ps {
	case PermStateGrant:
		return "GRANT"
	case PermStateGrantWithGrantOption:
		return "GRANT_WITH_GRANT_OPTION"
	case PermStateDeny:
		return "DENY"
	case PermStateRevoke:
		return "REVOKE"
	default:
		return ""
	}
}

// NewPermissionState конструктор PermissionState
func NewPermissionState(value string) PermissionState {
	switch value {
	case "GRANT":
		return PermStateGrant
	case "GRANT_WITH_GRANT_OPTION":
		return PermStateGrantWithGrantOption
	case "DENY":
		return PermStateDeny
	case "REVOKE":
		return PermStateRevoke
	default:
		return PermStateUnknown
	}
}

// PermStates состояния разрешений (GRANT, DENY etc)
type PermStates map[PermissionState]Permissions

// UserPerms разрешения пользователя
type UserPerms map[string]PermStates

// ObjectPermissions все разрешения на объект БД
type ObjectPermissions map[string]UserPerms

// Append добавляет информацию о разрешении
func (perms ObjectPermissions) Append(schema, object, permission, state, user string) error {
	obj := SchemaAndObject(schema, object, true)

	permState := NewPermissionState(state)

	if permState == PermStateUnknown {
		return fmt.Errorf("unknown permission state %s", state)
	}

	return perms.append(obj, permission, permState, user)
}

func (perms ObjectPermissions) append(objectName, permission string, state PermissionState, user string) error {
	if userPerms, ok := perms[objectName]; ok {
		if permStates, ok := userPerms[user]; ok {
			if permissions, ok := permStates[state]; ok {
				permissions[permission] = true
			} else {
				permissions := make(Permissions)
				permissions[permission] = true

				permStates[state] = permissions
			}
		} else {
			permissions := make(Permissions)
			permissions[permission] = true

			permStates := make(PermStates)
			permStates[state] = permissions

			userPerms[user] = permStates
		}
	} else {
		permissions := make(Permissions)
		permissions[permission] = true

		permStates := make(PermStates)
		permStates[state] = permissions

		userPerms := make(UserPerms)
		userPerms[user] = permStates

		perms[objectName] = userPerms
	}

	return nil
}

// Users возвращает список пользователей, обладающих правами на указанный объект
func (perms UserPerms) Users() []string {
	users := make([]string, len(perms))
	var index int

	for user := range perms {
		users[index] = user
		index++
	}

	return users
}

const selectPermissions = `
select permissions.[schema], permissions.object, permissions.permission, permissions.state, permissions.[user]
from (
    select
        [catalog] = db_name(),
        [schema] = iif(perm.class = 1, schema_name(objects.schema_id), null),
        [object] = case class
            when 0 then db_name()
            when 1 then object_name(objects.object_id)
            when 3 then schema_name(perm.major_id)
            else null
        end,
        [permission] = perm.permission_name,
        [state] = perm.state_desc,
        [user] = user_name(grantee_principal_id)
    from sys.database_permissions as perm
        left join sys.objects as objects on perm.major_id = objects.object_id
    where perm.major_id > 0
) as permissions
where not permissions.catalog is null and not permissions.object is null
    and not permissions.permission is null and not permissions.state is null
    and not permissions.[user] is null
order by permissions.catalog, permissions.[schema], permissions.object,
    permissions.[user], permissions.state, permissions.permission
`
