import { useSessionStore } from '@shared/session';
import { u } from 'framer-motion/client';
import React, { useEffect, useState } from 'react'
import { Link, useLocation } from 'react-router-dom';

class AppMenuItem {
    id: number | undefined;
    parentId: number | undefined;
    route: string;
    label: string;
    icon: string | undefined;
    permissionName: string | undefined;
    isActive: boolean = false;
    isCollapsed: boolean = false;
    children: AppMenuItem[];

    constructor(
        label: string,
        route: string,
        icon?: string,
        permissionName?: string,
        children: AppMenuItem[] = []
    ) {
        this.route = route;
        this.label = label;
        this.icon = icon;
        this.permissionName = permissionName;
        this.children = children;
    }
}

const defaultMenuItems: AppMenuItem[] = [
    new AppMenuItem('About', '/app/home', 'fas fa-info-circle'),
    new AppMenuItem('Categories', '/app/main/categories', 'fas fa-tags', 'Pages.Categories'),
    new AppMenuItem('Products', '/app/main/products', 'fas fa-store', 'Pages.Products'),
    new AppMenuItem('Administration', '', 'fas fa-tasks', '', [
        new AppMenuItem('Roles', '/app/admin/roles', 'fas fa-layer-group', 'Pages.Administration.Roles'),
        new AppMenuItem('Users', '/app/admin/users', 'fas fa-users', 'Pages.Administration.Users'),
    ])
];
const menuItemsMap: { [key: number]: AppMenuItem } = {};
let activatedMenuItems: AppMenuItem[] = [];

const patchMenuItems = (items: AppMenuItem[], parentId?: number) => {
    items.forEach((item, index) => {
        item.id = parentId ? Number(parentId + '' + (index + 1)) : index + 1;
        if (parentId) {
            item.parentId = parentId;
        }
        if (parentId || item.children.length > 0) {
            menuItemsMap[item.id] = item;
        }
        if (item.children.length > 0) {
            patchMenuItems(item.children, item.id);
        }
    });
}

const activateMenuItems = (url: string, items: AppMenuItem[]) => {
    deactivateMenuItems(items);
    activatedMenuItems = [];
    const foundedItems = findMenuItemsByUrl(url, items);
    foundedItems.forEach((item) => {
        activateMenuItem(item);
    });
}

const deactivateMenuItems = (items: AppMenuItem[]) => {
    items.forEach((item: AppMenuItem) => {
        item.isActive = false;
        item.isCollapsed = true;
        if (item.children.length > 0) {
            deactivateMenuItems(item.children);
        }
    });
}

const activateMenuItem = (item: AppMenuItem): void => {
    item.isActive = true;
    if (item.children) {
        item.isCollapsed = false;
    }
    activatedMenuItems.push(item);
    if (item.parentId) {
        activateMenuItem(menuItemsMap[item.parentId]);
    }
}

const findMenuItemsByUrl = (
    url: string,
    items: AppMenuItem[],
    foundedItems: AppMenuItem[] = []
) => {
    if (url == "/" || url == "/app") {
        url = "/app/home";
    }
    items.forEach((item: AppMenuItem) => {
        if (item.route != '' && url.indexOf(item.route) !== -1) {
            foundedItems.push(item);
        } else if (item.children) {
            findMenuItemsByUrl(url, item.children, foundedItems);
        }
    });

    return foundedItems;
}

const refreshMenuItems = (url: string, items: AppMenuItem[]) => {
    patchMenuItems(items);
    activateMenuItems(url, items);
}

const SidebarMenu = () => {
    const location = useLocation();
    const [menuItems, setMenuItems] = useState(() => {
        refreshMenuItems(location.pathname, defaultMenuItems);
        return defaultMenuItems;
    });
    const { isGranted } = useSessionStore();

    const handleCollapsed = (item: AppMenuItem) => {
        setMenuItems((prevMenuItems) => {
            const newMenuItems = prevMenuItems.map((menuItem) => {
                if (menuItem.id === item.id) {
                    return { ...menuItem, isCollapsed: !menuItem.isCollapsed };
                }
                return menuItem;
            });
            return newMenuItems;
        });
    };

    const checkIsMenuItemVisible = (item: AppMenuItem) => {
        if (!item.permissionName && (!item.children || item.children.length == 0)) {
            return true;
        }

        if (!item.permissionName && item.children && item.children.length > 0) {
            return item.children.some(x => isGranted(x.permissionName ?? ""));
        }

        return isGranted(item.permissionName ?? "");
    }

    useEffect(() => {
        setMenuItems((prevMenuItems: AppMenuItem[]) => {
            const newMenuItems = [...prevMenuItems];
            refreshMenuItems(location.pathname, newMenuItems);
            return newMenuItems;
        });
    }, [location]);

    return (
        <>
            {
                menuItems.map((item) =>
                    checkIsMenuItemVisible(item) &&
                    <div key={'MenuItem' + item.id} className={`menu-item menu-accordion${item.isActive ? ' here' : ''}${!item.isCollapsed ? ' show' : ''}`}>
                        {
                            item.children.length == 0 &&
                            <Link to={item.route} className={`menu-link${item.isActive ? ' active' : ''}`}>
                                <span className="menu-icon">
                                    <i className={item.icon}></i>
                                </span>
                                <span className="menu-title">{item.label}</span>
                            </Link>
                        }
                        {
                            item.children.length > 0 &&
                            <span className={`menu-link${item.isActive ? ' active' : ''}`} onClick={() => handleCollapsed(item)}>
                                <span className="menu-icon">
                                    <i className={item.icon}></i>
                                </span>
                                <span className="menu-title">{item.label}</span>
                                <span className="menu-arrow"></span>
                            </span>
                        }
                        {
                            item.children.length > 0 &&
                            <div className={`menu-sub menu-sub-accordion${!item.isCollapsed ? ' show' : ''}`}>
                                {
                                    item.children.map((childItem) =>
                                        checkIsMenuItemVisible(childItem) &&
                                        <div key={'SubMenuItem' + childItem.id} className="menu-item">
                                            <Link to={childItem.route} className={`menu-link${childItem.isActive ? ' active' : ''}`}>
                                                <span className="menu-icon">
                                                    <i className={childItem.icon}></i>
                                                </span>
                                                <span className="menu-title">{childItem.label}</span>
                                            </Link>
                                        </div>
                                    )
                                }
                            </div>
                        }
                    </div>
                )
            }
        </>
    )
}

export default SidebarMenu