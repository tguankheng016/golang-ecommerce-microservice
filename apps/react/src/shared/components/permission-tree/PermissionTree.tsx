import { APIClient } from "@shared/service-proxies";
import { PermissionGroupDto } from "@shared/service-proxies/identity-service-proxies";
import { Dispatch, RefObject, SetStateAction, useEffect, useRef, useState } from "react";
import styles from './PermissionTree.module.css';
import { Link } from "react-router-dom";
import { TreeNode } from "primereact/treenode";
import { Tree, TreeCheckboxSelectionKeys, TreeEventNodeEvent } from "primereact/tree";

class ExtendedPermissionGroupDto extends PermissionGroupDto {
    isActive: boolean = false;
    selectedCount: number = 0;
}

interface Props {
    show: boolean;
    loading: boolean;
    grantedPermissions: string[];
    setGrantedPermissions: Dispatch<SetStateAction<string[]>>;
}

const PermissionTree = ({ show, loading, grantedPermissions, setGrantedPermissions }: Props) => {
    const isGrantAllCheckboxRef = useRef<HTMLInputElement>(null);
    const isSelectAllCheckboxRef = useRef<HTMLInputElement>(null);

    const [allPermissions, setAllPermissions] = useState<ExtendedPermissionGroupDto[]>([]);
    const [permissionNodes, setPermissionNodes] = useState<TreeNode[]>([]);
    const [selectedKeys, setSelectedKeys] = useState<TreeCheckboxSelectionKeys | null>(null);
    const [isGrantAll, setIsGrantAll] = useState<boolean>(false);
    const [isSelectAll, setIsSelectAll] = useState<boolean>(false);
    const [activePermissionGroup, setActivePermissionGroup] = useState<ExtendedPermissionGroupDto>();
    const [allPermissionsCount, setAllPermissionsCount] = useState<number>(0);

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        if (!show || loading) {
            reset();
            return;
        }

        const identityService = APIClient.getIdentityService();
        identityService.getAllPermissions(signal).then((response) => {
            if (response.items) {
                const extendedPermissions = response.items.map((item: PermissionGroupDto) => {
                    const extendedItem = new ExtendedPermissionGroupDto();
                    Object.assign(extendedItem, item);
                    return extendedItem;
                });
                setAllPermissions(extendedPermissions);

                if (extendedPermissions.length > 0) {
                    const firstGroup = extendedPermissions[0];
                    setActive(firstGroup);
                    const extendedPermissionsCount = extendedPermissions.map(x => x.permissions).flat().length;
                    setAllPermissionsCount(extendedPermissionsCount);

                    if (grantedPermissions.length > 0) {
                        if (grantedPermissions.length == extendedPermissionsCount) {
                            setIsGrantAll(true);
                        } else {
                            setCheckboxIntermediate(isGrantAllCheckboxRef);
                        }
                    }

                    setAllPermissions((prevPermissions) => {
                        return prevPermissions.map((p) => {
                            if (p.groupName !== firstGroup.groupName) {
                                p.selectedCount = 0;

                                p.permissions?.forEach((i) => {
                                    p.selectedCount += (i.name && grantedPermissions.indexOf(i.name) !== -1 ? 1 : 0);
                                });
                            }
                            return p;
                        });
                    });
                }
            }
        });

        return () => {
            abortController.abort();
        };
    }, [show, loading]);

    const setActive = (group: ExtendedPermissionGroupDto) => {
        setAllPermissions((prevPermissions) => {
            return prevPermissions.map((p) => {
                if (activePermissionGroup && p.groupName === activePermissionGroup.groupName) {
                    p.selectedCount = activePermissionGroup.selectedCount;
                }
                p.isActive = p.groupName === group.groupName;
                return p;
            });
        });

        populateTreeNodes(group);
    };

    const populateTreeNodes = (group: ExtendedPermissionGroupDto) => {
        const nodes: TreeNode[] = [];
        let selectedNodes: TreeNode[] = [];

        group.selectedCount = 0;

        group.permissions?.forEach((p) => {
            const newNode = {
                key: p.name,
                label: p.displayName,
                data: p.name,
            };

            group.selectedCount += p.name && grantedPermissions.indexOf(p.name) !== -1 ? 1 : 0;
            nodes.push(newNode);
        });

        selectedNodes = nodes.filter((p) => grantedPermissions.indexOf(p.data) !== -1);

        setPermissionNodes(nodes);
        setSelectedKeys(extractKeys(selectedNodes));
        setActivePermissionGroup(group);
        populateSelectAllCheckboxStatus(group);
    }

    const populateSelectAllCheckboxStatus = (group: ExtendedPermissionGroupDto) => {
        setIsSelectAll(false);
        setCheckboxIntermediate(isSelectAllCheckboxRef, false);

        if (group.selectedCount == group.permissions?.length) {
            setIsSelectAll(true);
        } else if (group.selectedCount && group.selectedCount > 0) {
            setCheckboxIntermediate(isSelectAllCheckboxRef);
        }
    }

    const extractKeys = (nodes: TreeNode[]) => {
        const keys: TreeCheckboxSelectionKeys = {};
        nodes.forEach((node) => {
            if (node.key) {
                keys[node.key] = {
                    checked: true,
                    partialChecked: false
                };
            }
        });

        return keys;
    }

    const onNodeSelect = (e: TreeEventNodeEvent) => {
        if (e?.node?.data) {
            setActivePermissionGroup(prevState => {
                if (!prevState) return prevState;
                const newState = new ExtendedPermissionGroupDto();
                Object.assign(newState, prevState);
                newState.selectedCount += 1;
                return newState;
            });

            setGrantedPermissions(prevState => [...prevState, e.node.data as string]);

            if ((activePermissionGroup?.selectedCount ?? 0) + 1 === activePermissionGroup?.permissions?.length) {
                setIsSelectAll(true);
                setCheckboxIntermediate(isSelectAllCheckboxRef, false);

                if (grantedPermissions.length + 1 === allPermissionsCount) {
                    setIsGrantAll(true);
                    setCheckboxIntermediate(isGrantAllCheckboxRef, false);
                }
            } else {
                setCheckboxIntermediate(isSelectAllCheckboxRef);
                setCheckboxIntermediate(isGrantAllCheckboxRef);
            }
        }
    }

    const onNodeUnselect = (e: TreeEventNodeEvent) => {
        if (e?.node?.data) {
            setActivePermissionGroup(prevState => {
                if (!prevState) return prevState;
                const newState = new ExtendedPermissionGroupDto();
                Object.assign(newState, prevState);
                newState.selectedCount -= 1;
                return newState;
            });

            setGrantedPermissions(prevState => prevState.filter(p => p !== e.node.data));

            if ((activePermissionGroup?.selectedCount ?? 0) - 1 === 0) {
                setIsSelectAll(false);
                setCheckboxIntermediate(isSelectAllCheckboxRef, false);

                if (grantedPermissions.length - 1 === 0) {
                    setIsGrantAll(false);
                    setCheckboxIntermediate(isGrantAllCheckboxRef, false);
                }
            } else {
                setCheckboxIntermediate(isSelectAllCheckboxRef);
                setCheckboxIntermediate(isGrantAllCheckboxRef);
            }
        }
    }

    const onSelectAll = (checked: boolean) => {
        setIsSelectAll(checked);

        if (checked) {
            const selectedNodes = [...permissionNodes];
            setSelectedKeys(extractKeys(selectedNodes));
            setGrantedPermissions(prevState => {
                const newPermissions = [...prevState];
                selectedNodes.forEach((n) => {
                    if (n.data && newPermissions.indexOf(n.data as string) === -1) {
                        newPermissions.push(n.data as string);
                    }
                });
                return newPermissions;
            });
            setActivePermissionGroup(prevState => {
                if (!prevState) return prevState;
                const newState = new ExtendedPermissionGroupDto();
                Object.assign(newState, prevState);
                newState.selectedCount = newState.permissions?.length ?? 0;
                return newState;
            });

            if (grantedPermissions.length + selectedNodes.length === allPermissionsCount) {
                setIsGrantAll(true);
                setCheckboxIntermediate(isGrantAllCheckboxRef, false);
            } else {
                setCheckboxIntermediate(isGrantAllCheckboxRef);
            }
        } else {
            setSelectedKeys(null);
            setGrantedPermissions(prevState => prevState.filter(p => permissionNodes.map(n => n.data as string).indexOf(p) === -1));
            setActivePermissionGroup(prevState => {
                if (!prevState) return prevState;
                const newState = new ExtendedPermissionGroupDto();
                Object.assign(newState, prevState);
                newState.selectedCount = 0;
                return newState;
            });
            setIsGrantAll(false);

            if (grantedPermissions.length == permissionNodes.length) {
                setCheckboxIntermediate(isGrantAllCheckboxRef, false);
            } else {
                setCheckboxIntermediate(isGrantAllCheckboxRef);
            }
        }
    }

    const onGrantAll = (checked: boolean) => {
        setIsGrantAll(checked);
        setCheckboxIntermediate(isGrantAllCheckboxRef, false);
        setCheckboxIntermediate(isSelectAllCheckboxRef, false);

        if (checked) {
            onSelectAll(true);
            setAllPermissions((prevPermissions) => {
                return prevPermissions.map((p) => {
                    if (activePermissionGroup && p.groupName !== activePermissionGroup.groupName) {
                        p.selectedCount = p.permissions?.length ?? 0;
                    }
                    return p;
                });
            });
            setGrantedPermissions([...allPermissions.map(p => p.permissions).flat().map(p => p?.name ?? "")]);
            setCheckboxIntermediate(isGrantAllCheckboxRef, false);
        } else {
            onSelectAll(false);
            setAllPermissions((prevPermissions) => {
                return prevPermissions.map((p) => {
                    if (activePermissionGroup && p.groupName !== activePermissionGroup.groupName) {
                        p.selectedCount = 0;
                    }
                    return p;
                });
            });
            setGrantedPermissions([]);
            setCheckboxIntermediate(isGrantAllCheckboxRef, false);
        }
    }

    const setCheckboxIntermediate = (checkbox: RefObject<HTMLInputElement>, setIndeterminate: boolean = true): void => {
        if (checkbox && checkbox.current)
            checkbox.current.indeterminate = setIndeterminate;
    };

    const reset = () => {
        setIsGrantAll(false);
        setIsSelectAll(false);
        setSelectedKeys(null);
        setAllPermissions([]);
        setActivePermissionGroup(undefined);
    }

    return (
        <>
            <div className="row">
                <div className="col-4">
                    <label className="form-check form-check-custom form-check-solid py-2">
                        <input
                            ref={isGrantAllCheckboxRef}
                            type="checkbox"
                            name="IsGrantAllPermissions"
                            className="form-check-input"
                            checked={isGrantAll}
                            onChange={(e) => onGrantAll(e.target.checked)}
                        />
                        <span className="fw-semibold ps-2 fs-6">
                            Grant all permissions
                        </span>
                    </label>
                    <hr />
                </div>
                <div className="col-8">
                    <label className="form-check form-check-custom form-check-solid py-2">
                        <input
                            ref={isSelectAllCheckboxRef}
                            type="checkbox"
                            name="SelectAllPermissionsUnderGroup"
                            className="form-check-input"
                            checked={isSelectAll}
                            onChange={(e) => onSelectAll(e.target.checked)}
                        />
                        <span className="fw-semibold ps-2 fs-6">
                            Select all
                        </span>
                    </label>
                    <hr />
                </div>
            </div>
            <div className="row">
                <div className="col-4">
                    {allPermissions.map((group) => (
                        <div key={group.groupName} className={`menu-item menu-accordion ${styles["menu-item"]}`}>
                            <Link
                                to="#"
                                className={`menu-link ${styles["menu-link"]} ${group.isActive ? styles.active : ''}`}
                                onClick={() => setActive(group)}
                            >
                                <span className="text-truncate">
                                    {group.groupName}
                                    {
                                        group.isActive && activePermissionGroup ?
                                            (activePermissionGroup.selectedCount > 0 && <span> ({activePermissionGroup.selectedCount})</span>) :
                                            (group.selectedCount > 0 && <span> ({group.selectedCount})</span>)
                                    }
                                </span>
                            </Link>
                        </div>
                    ))}
                </div>
                <div className="col-8">
                    {
                        permissionNodes && permissionNodes.length > 0 &&
                        <Tree
                            value={permissionNodes}
                            selectionMode="checkbox"
                            selectionKeys={selectedKeys}
                            onSelectionChange={e => setSelectedKeys(e.value as TreeCheckboxSelectionKeys | null)}
                            onSelect={(e) => onNodeSelect(e)}
                            onUnselect={onNodeUnselect}
                            propagateSelectionUp={false}
                        />
                    }
                </div>
            </div>
        </>
    )
}

export default PermissionTree