import { DefaultPage, BreadcrumbItem } from '@app/components/layout';
import { AdvancedFilter, AdvancedFilterProps } from '@shared/components/advanced-filter';
import { useDataTable, PrimengTableHelper, DateTimeBodyTemplate } from '@shared/primeng';
import APIClient from '@shared/service-proxies/api-client';
import { RoleDto } from '@shared/service-proxies/identity-service-proxies';
import { SwalMessageService, SwalNotifyService } from '@shared/sweetalert2';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { MenuItem } from 'primereact/menuitem';
import { Paginator } from 'primereact/paginator';
import { TieredMenu } from 'primereact/tieredmenu';
import { useEffect, useRef, useState } from 'react';
import CreateOrEditRoleModal from './CreateOrEditRoleModal';
import { useSessionStore } from '@shared/session';
import { Skeleton } from 'primereact/skeleton';

interface RoleTableProps {
    filterText: string | undefined;
    reloading: boolean;
    getMenuItems: (item: RoleDto) => MenuItem[];
}

const RoleAdvancedFilter = ({ filterText, setFilterText }: AdvancedFilterProps) => {
    return (
        <AdvancedFilter
            filterText={filterText}
            setFilterText={setFilterText}
        >
        </AdvancedFilter>
    )
}

const RoleTable = ({ filterText, reloading, getMenuItems }: RoleTableProps) => {
    const menu = useRef<TieredMenu>(null);
    const [menuItems, setMenuItems] = useState<MenuItem[]>([]);
    const {
        loading,
        setLoading,
        totalRecords,
        setTotalRecords,
        lazyState,
        onSort,
        onPageChange,
    } = useDataTable();
    const [roles, setRoles] = useState<RoleDto[]>([]);

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        loadLazyData(signal);

        return () => {
            abortController.abort();
        };
    }, [lazyState, filterText, reloading]);

    const loadLazyData = (signal: AbortSignal) => {
        const roleService = APIClient.getRoleService();

        const loadingTimer = setTimeout(() => {
            setLoading(true);
        }, 200);

        roleService.getRoles(
            lazyState.rows,
            lazyState.first,
            PrimengTableHelper.getSorting(lazyState),
            filterText,
            signal
        ).then((res) => {
            setRoles(res.items ?? []);
            setTotalRecords(res.totalCount ?? 0);
        }).finally(() => {
            clearTimeout(loadingTimer);
            setLoading(false);
        });
    };

    const handleButtonClick = (event: React.MouseEvent<HTMLButtonElement>, record: RoleDto) => {
        setMenuItems(getMenuItems(record));
        if (menu.current) {
            menu.current.toggle(event);
        }
    };

    const actionButtonBodyTemplate = (rowData: RoleDto) => {
        if (loading) {
            return <Skeleton></Skeleton>;
        }

        return (
            <div className="btn-group dropdown">
                <button className="dropdown-toggle btn btn-sm btn-primary" onClick={(e) => handleButtonClick(e, rowData)}>
                    <i className="fa fa-cog"></i>
                    <span className="caret"></span>
                    Actions
                </button>
                <TieredMenu model={menuItems} popup ref={menu} appendTo={document.body} />
            </div>
        );
    }

    const roleNameTemplate = (rowData: RoleDto) => {
        if (loading) {
            return <Skeleton></Skeleton>;
        }

        return (
            <>
                <span className="p-column-title">Role name</span>
                <span>{rowData.name}</span>
                {
                    rowData.isStatic &&
                    <span className="badge badge-light-success fw-bold fs-8 px-2 py-1 ms-2">Static</span>
                }
                {
                    rowData.isDefault &&
                    <span className="badge badge-light-primary fw-bold fs-8 px-2 py-1 ms-2">Default</span>
                }
            </>
        );
    }

    return (
        <div className="row align-items-center mx-0">
            <div className="primeng-datatable-container">
                <DataTable
                    value={roles}
                    lazy
                    rows={lazyState.rows}
                    onSort={onSort}
                    sortField={lazyState.sortField}
                    sortOrder={lazyState.sortOrder}
                >
                    <Column header="Actions" body={actionButtonBodyTemplate} style={{ width: '130px' }} />
                    <Column field="name" header="Role Name" sortable body={roleNameTemplate} />
                    <Column field="created_at" header="Creation Time" sortable body={(data: RoleDto) => DateTimeBodyTemplate(data.createdAt, "Creation Time", loading)} />
                </DataTable>
                {
                    totalRecords == 0 && !loading &&
                    <div className="primeng-no-data">
                        No record found
                    </div>
                }
                <Paginator
                    template={PrimengTableHelper.defaultPaginatorTemplate}
                    first={lazyState.first}
                    rows={lazyState.rows}
                    totalRecords={totalRecords}
                    rowsPerPageOptions={PrimengTableHelper.predefinedRecordsCountPerPage}
                    onPageChange={onPageChange}
                    currentPageReportTemplate={PrimengTableHelper.defaultCurrentReportTemplate}
                />
            </div>
        </div>
    )
}

const RolePage = () => {
    const [filterText, setFilterText] = useState('');
    const [reloading, setReloading] = useState(false);
    const [showModal, setShowModal] = useState(false);
    const [roleId, setRoleId] = useState(0);
    const { isGranted } = useSessionStore();

    const breadcrumbs: BreadcrumbItem[] = [
        new BreadcrumbItem('Administration'),
        new BreadcrumbItem('Roles')
    ];

    const getMenuItemsForItem = (item: RoleDto): MenuItem[] => {
        return [
            {
                label: 'Edit',
                command: (event) => {
                    setRoleId(item.id ?? 0);
                    setShowModal(true);
                },
                visible: isGranted('Pages.Administration.Roles.Edit')
            },
            {
                label: 'Delete',
                command: (event) => {
                    handleDelete(item);
                },
                visible: isGranted('Pages.Administration.Roles.Delete')
            }
        ];
    };

    const actionButtons = () => {
        if (!isGranted("Pages.Administration.Roles.Create")) {
            return null;
        }

        return (
            <div>
                <button
                    className="btn btn-sm btn-primary"
                    onClick={handleCreate}
                >
                    <i className="fa fa-plus btn-md-icon"></i>
                    <span className="d-none d-md-inline-block">Create</span>
                </button>
            </div>
        )
    }

    const handleCreate = () => {
        setRoleId(0);
        setShowModal(true);
    }

    const handleDelete = (role: RoleDto) => {
        SwalMessageService.showConfirmation("Are you sure?", `Role "${role.name}" will be deleted`, () => {
            const roleService = APIClient.getRoleService();
            roleService.deleteRole(role.id ?? 0).then(() => {
                SwalNotifyService.success('Deleted successfully');
                setReloading(!reloading);
            });
        });
    }

    return (
        <>
            <DefaultPage title="Roles" breadcrumbs={breadcrumbs} actionButtons={actionButtons()}>
                <RoleAdvancedFilter filterText={filterText} setFilterText={setFilterText} />
                <RoleTable reloading={reloading} filterText={filterText} getMenuItems={getMenuItemsForItem} />
            </DefaultPage>
            <CreateOrEditRoleModal roleId={roleId} show={showModal} handleClose={() => setShowModal(false)} handleSave={() => setReloading(!reloading)} />
        </>
    )
}

export default RolePage