import { u } from 'framer-motion/client';
import React, { useState } from 'react'
import { Link } from 'react-router-dom';

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
    new AppMenuItem('Administration', '', 'fas fa-tasks', '', [
        new AppMenuItem('Roles', '/app/admin/roles', 'fas fa-layer-group', 'Pages.Administration.Roles'),
        new AppMenuItem('Users', '/app/admin/users', 'fas fa-users','Pages.Administration.Users'),
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

const activateMenuItems = (url: string) => {
    deactivateMenuItems(defaultMenuItems);
    activatedMenuItems = [];
    const foundedItems = findMenuItemsByUrl(url, defaultMenuItems);
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

const SidebarMenu = () => {
    const [menuItems, setMenuItems] = useState(() => {
        patchMenuItems(defaultMenuItems);
        activateMenuItems(location.pathname);

        return defaultMenuItems;
    });

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

    return (
        <>
            {
                menuItems.map((item) => 
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
                                    <i className={ item.icon }></i>
                                </span>
                                <span className="menu-title">{ item.label }</span>
                                <span className="menu-arrow"></span>
                            </span>
                        }
                        {
                            item.children.length > 0 && 
                            <div className={`menu-sub menu-sub-accordion${!item.isCollapsed ? ' show' : ''}`}>
                                {
                                    item.children.map((childItem) => 
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
    // return (
    //     <div>
    //         <div 
    //             className="menu menu-rounded menu-column menu-md-row menu-state-bg menu-title-white menu-state-icon-primary menu-state-bullet-primary menu-arrow-gray-400 fw-bold px-4 px-lg-0 my-5 my-lg-0 align-items-stretch" 
    //             id="#kt_header_menu" 
    //             data-kt-menu="true"
    //         >
    //             {
    //                 menuItems.map((item) => 
    //                     <Link key={item.label + 'MenuItem'} to={item.route} className={ checkIsActive(item) ? activeMenuItemClassName : menuItemClassName}>
    //                         <span className="menu-link py-3">
    //                             <span className="menu-title">{item.label}</span>
    //                         </span>
    //                     </Link>
    //                 )
    //             }
    //         </div>
    //     </div>
    // )
}

export default SidebarMenu