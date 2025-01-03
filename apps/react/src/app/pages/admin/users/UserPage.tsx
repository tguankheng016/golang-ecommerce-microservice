import { DefaultPage, BreadcrumbItem } from '@app/components/layout';
import { TextBodyTemplate } from '@shared/primeng';
import { useDataTable, PrimengTableHelper } from '@shared/primeng';
import APIClient from '@shared/service-proxies/api-client';
import { UserDto } from '@shared/service-proxies/identity-service-proxies';
import StringHelper from '@shared/utils/string-helper';
import { Column } from 'primereact/column';
import { DataTable } from 'primereact/datatable';
import { MenuItem } from 'primereact/menuitem';
import { Paginator } from 'primereact/paginator';
import { TieredMenu } from 'primereact/tieredmenu';
import { useEffect, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import CreateOrEditUserModal from './CreateOrEditUserModal';
import EditUserPermissionsModal from './EditUserPermissionsModal';
import { SwalMessageService, SwalNotifyService } from '@shared/sweetalert2';
import { AdvancedFilter, AdvancedFilterProps } from '@shared/components/advanced-filter';
import { useSessionStore } from '@shared/session';
import { Skeleton } from 'primereact/skeleton';

interface UserTableProps {
    filterText: string | undefined;
    reloading: boolean;
    getMenuItems: (item: UserDto) => MenuItem[];
}

const UserAdvancedFilter = ({ filterText, setFilterText }: AdvancedFilterProps) => {
    return (
        <AdvancedFilter
            filterText={filterText}
            setFilterText={setFilterText}
        >
        </AdvancedFilter>
    )
}

const UserTable = ({ filterText, reloading, getMenuItems }: UserTableProps) => {
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
    const [users, setUsers] = useState<UserDto[]>([]);

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        loadLazyData(signal);

        return () => {
            abortController.abort();
        };
    }, [lazyState, filterText, reloading]);

    const loadLazyData = (signal: AbortSignal) => {
        const userService = APIClient.getUserService();

        const loadingTimer = setTimeout(() => {
            setLoading(true);
        }, 200);

        userService.getUsers(
            lazyState.rows,
            lazyState.first,
            PrimengTableHelper.getSorting(lazyState),
            filterText,
            signal
        ).then((res) => {
            setUsers(res.items ?? []);
            setTotalRecords(res.totalCount ?? 0);
        }).finally(() => {
            clearTimeout(loadingTimer);
            setLoading(false);
        });
    };

    const handleButtonClick = (event: React.MouseEvent<HTMLButtonElement>, record: UserDto) => {
        setMenuItems(getMenuItems(record));
        if (menu.current) {
            menu.current.toggle(event);
        }
    };

    const actionButtonBodyTemplate = (rowData: UserDto) => {
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

    const usernameBodyTemplate = (rowData: UserDto) => {
        if (loading) {
            return <Skeleton></Skeleton>;
        }

        const profileImageUrl = setUsersProfilePictureUrl(rowData);
        return (
            <>
                <span className="p-column-title">Username</span>
                <div className="w-md-100 d-flex align-items-center">
                    <Link to={profileImageUrl} target='_blank'>
                        <img className="img-circle rounded-circle" src={profileImageUrl} />
                    </Link>
                    <span>
                        {rowData.userName}
                    </span >
                </div >
            </>
        );
    };

    const setUsersProfilePictureUrl = (user: UserDto) => {
        return StringHelper.formatString(
            import.meta.env.VITE_UI_AVATAR_URL,
            user?.firstName ?? "",
            user?.lastName ?? ""
        );
    }

    return (
        <div className="row align-items-center mx-0">
            <div className="primeng-datatable-container">
                <DataTable
                    value={users}
                    lazy
                    rows={lazyState.rows}
                    onSort={onSort}
                    sortField={lazyState.sortField}
                    sortOrder={lazyState.sortOrder}
                >
                    <Column header="Actions" body={actionButtonBodyTemplate} style={{ width: '130px' }} />
                    <Column field="user_name" header="Username" sortable body={usernameBodyTemplate} />
                    <Column field="first_name" header="First name" sortable body={(data: UserDto) => TextBodyTemplate(data.firstName, "First name", loading)} />
                    <Column field="last_name" header="Last name" sortable body={(data: UserDto) => TextBodyTemplate(data.lastName, "Last name", loading)} />
                    <Column field="email" header="Email" sortable body={(data: UserDto) => TextBodyTemplate(data.email, "Email", loading)} />
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

const UserPage = () => {
    const [filterText, setFilterText] = useState('');
    const [showModal, setShowModal] = useState(false);
    const [showPermissionModal, setShowPermissionModal] = useState(false);
    const [userId, setUserId] = useState(0);
    const [userDto, setUserDto] = useState(new UserDto());
    const [reloading, setReloading] = useState(false);
    const { isGranted } = useSessionStore();

    const breadcrumbs: BreadcrumbItem[] = [
        new BreadcrumbItem('Administration'),
        new BreadcrumbItem('Users')
    ];

    const getMenuItemsForItem = (item: UserDto): MenuItem[] => {
        return [
            {
                label: 'Edit',
                command: (event) => {
                    setUserId(item.id ?? 0);
                    setShowModal(true);
                },
                visible: isGranted('Pages.Administration.Users.Edit')
            },
            {
                label: 'Permissions',
                command: (event) => {
                    setUserDto(item);
                    setShowPermissionModal(true);
                },
                visible: isGranted('Pages.Administration.Users.ChangePermissions')
            },
            {
                label: 'Delete',
                command: (event) => {
                    handleDelete(item);
                },
                visible: isGranted('Pages.Administration.Users.Delete')
            }
        ];
    };

    const actionButtons = () => {
        if (!isGranted('Pages.Administration.Users.Create')) {
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
        setUserId(0);
        setShowModal(true);
    }

    const handleDelete = (user: UserDto) => {
        SwalMessageService.showConfirmation("Are you sure?", `User "${user.userName}" will be deleted`, () => {
            const userService = APIClient.getUserService();
            userService.deleteUser(user.id ?? 0).then(() => {
                SwalNotifyService.success('Deleted successfully');
                setReloading(!reloading);
            });
        });
    }

    return (
        <>
            <DefaultPage title="Users" breadcrumbs={breadcrumbs} actionButtons={actionButtons()}>
                <UserAdvancedFilter filterText={filterText} setFilterText={setFilterText} />
                <UserTable reloading={reloading} filterText={filterText} getMenuItems={getMenuItemsForItem} />
            </DefaultPage>
            <CreateOrEditUserModal userId={userId} show={showModal} handleClose={() => setShowModal(false)} handleSave={() => setReloading(!reloading)} />
            <EditUserPermissionsModal show={showPermissionModal} handleClose={() => setShowPermissionModal(false)} userDto={userDto} />
        </>
    )
}

export default UserPage