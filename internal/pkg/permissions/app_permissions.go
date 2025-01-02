package permissions

const (
	UserPermissionsGroupName                  = "Users"
	PagesAdministrationUsers                  = "Pages.Administration.Users"
	PagesAdministrationUsersCreate            = "Pages.Administration.Users.Create"
	PagesAdministrationUsersEdit              = "Pages.Administration.Users.Edit"
	PagesAdministrationUsersDelete            = "Pages.Administration.Users.Delete"
	PagesAdministrationUsersChangePermissions = "Pages.Administration.Users.ChangePermissions"

	RolePermissionsGroupName       = "Roles"
	PagesAdministrationRoles       = "Pages.Administration.Roles"
	PagesAdministrationRolesCreate = "Pages.Administration.Roles.Create"
	PagesAdministrationRolesEdit   = "Pages.Administration.Roles.Edit"
	PagesAdministrationRolesDelete = "Pages.Administration.Roles.Delete"

	CategoryPermissionsGroupName = "Categories"
	PagesCategories              = "Pages.Categories"
	PagesCategoriesCreate        = "Pages.Categories.Create"
	PagesCategoriesEdit          = "Pages.Categories.Edit"
	PagesCategoriesDelete        = "Pages.Categories.Delete"

	ProductPermissionsGroupName = "Products"
	PagesProducts               = "Pages.Products"
	PagesProductsCreate         = "Pages.Products.Create"
	PagesProductsEdit           = "Pages.Products.Edit"
	PagesProductsDelete         = "Pages.Products.Delete"
)

type Permission struct {
	Name        string
	DisplayName string
	Group       string
}

type AppPermissions struct {
	Items map[string]Permission
}

var permissions = map[string]Permission{
	// Users
	PagesAdministrationUsers: {
		Name:        PagesAdministrationUsers,
		DisplayName: "View Users",
		Group:       UserPermissionsGroupName,
	},
	PagesAdministrationUsersCreate: {
		Name:        PagesAdministrationUsersCreate,
		DisplayName: "Create Users",
		Group:       UserPermissionsGroupName,
	},
	PagesAdministrationUsersEdit: {
		Name:        PagesAdministrationUsersEdit,
		DisplayName: "Edit Users",
		Group:       UserPermissionsGroupName,
	},
	PagesAdministrationUsersDelete: {
		Name:        PagesAdministrationUsersDelete,
		DisplayName: "Delete Users",
		Group:       UserPermissionsGroupName,
	},
	PagesAdministrationUsersChangePermissions: {
		Name:        PagesAdministrationUsersChangePermissions,
		DisplayName: "Change User Permissions",
		Group:       UserPermissionsGroupName,
	},
	// Roles
	PagesAdministrationRoles: {
		Name:        PagesAdministrationRoles,
		DisplayName: "View Roles",
		Group:       RolePermissionsGroupName,
	},
	PagesAdministrationRolesCreate: {
		Name:        PagesAdministrationRolesCreate,
		DisplayName: "Create Roles",
		Group:       RolePermissionsGroupName,
	},
	PagesAdministrationRolesEdit: {
		Name:        PagesAdministrationRolesEdit,
		DisplayName: "Edit Roles",
		Group:       RolePermissionsGroupName,
	},
	PagesAdministrationRolesDelete: {
		Name:        PagesAdministrationRolesDelete,
		DisplayName: "Delete Roles",
		Group:       RolePermissionsGroupName,
	},
	// Categories
	PagesCategories: {
		Name:        PagesCategories,
		DisplayName: "View Categories",
		Group:       CategoryPermissionsGroupName,
	},
	PagesCategoriesCreate: {
		Name:        PagesCategoriesCreate,
		DisplayName: "Create Categories",
		Group:       CategoryPermissionsGroupName,
	},
	PagesCategoriesEdit: {
		Name:        PagesCategoriesEdit,
		DisplayName: "Edit Categories",
		Group:       CategoryPermissionsGroupName,
	},
	PagesCategoriesDelete: {
		Name:        PagesCategoriesDelete,
		DisplayName: "Delete Categories",
		Group:       CategoryPermissionsGroupName,
	},
	// Products
	PagesProducts: {
		Name:        PagesProducts,
		DisplayName: "View Products",
		Group:       ProductPermissionsGroupName,
	},
	PagesProductsCreate: {
		Name:        PagesProductsCreate,
		DisplayName: "Create Products",
		Group:       ProductPermissionsGroupName,
	},
	PagesProductsEdit: {
		Name:        PagesProductsEdit,
		DisplayName: "Edit Products",
		Group:       ProductPermissionsGroupName,
	},
	PagesProductsDelete: {
		Name:        PagesProductsDelete,
		DisplayName: "Delete Products",
		Group:       ProductPermissionsGroupName,
	},
}

func GetAppPermissions() AppPermissions {
	// Immutable
	return AppPermissions{Items: permissions}
}
