import { BreadcrumbItem, DefaultPage } from "@app/components/layout";
import { AdvancedFilter, AdvancedFilterProps } from "@shared/components/advanced-filter";
import { PrimengTableHelper, TextBodyTemplate, useDataTable } from "@shared/primeng";
import APIClient from "@shared/service-proxies/api-client";
import { CategoryDto } from "@shared/service-proxies/product-service-proxies";
import { useSessionStore } from "@shared/session";
import { SwalMessageService, SwalNotifyService } from "@shared/sweetalert2";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { MenuItem } from "primereact/menuitem";
import { Paginator } from "primereact/paginator";
import { TieredMenu } from "primereact/tieredmenu";
import { useEffect, useRef, useState } from "react";
import CreateOrEditCategoryModal from "./CreateOrEditCategoryModal";
import { Skeleton } from "primereact/skeleton";

interface CategoryTableProps {
    filterText: string | undefined;
    reloading: boolean;
    getMenuItems: (item: CategoryDto) => MenuItem[];
}

const CategoryAdvancedFilter = ({ filterText, setFilterText }: AdvancedFilterProps) => {
    return (
        <AdvancedFilter
            filterText={filterText}
            setFilterText={setFilterText}
        >
        </AdvancedFilter>
    )
}

const CategoryTable = ({ filterText, reloading, getMenuItems }: CategoryTableProps) => {
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
    const [categories, setCategories] = useState<CategoryDto[]>([]);

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        loadLazyData(signal);

        return () => {
            abortController.abort();
        };
    }, [lazyState, filterText, reloading]);

    const loadLazyData = (signal: AbortSignal) => {
        const categoryService = APIClient.getCategoryService();

        const loadingTimer = setTimeout(() => {
            setLoading(true);
        }, 200);

        categoryService.getCategories(
            lazyState.rows,
            lazyState.first,
            PrimengTableHelper.getSorting(lazyState),
            filterText,
            signal
        ).then((res) => {
            setCategories(res.items ?? []);
            setTotalRecords(res.totalCount ?? 0);
        }).finally(() => {
            clearTimeout(loadingTimer);
            setLoading(false);
        });
    };

    const handleButtonClick = (event: React.MouseEvent<HTMLButtonElement>, record: CategoryDto) => {
        setMenuItems(getMenuItems(record));
        if (menu.current) {
            menu.current.toggle(event);
        }
    };

    const actionButtonBodyTemplate = (rowData: CategoryDto) => {
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

    return (
        <div className="row align-items-center mx-0">
            <div className="primeng-datatable-container">
                <DataTable
                    value={categories}
                    lazy
                    rows={lazyState.rows}
                    onSort={onSort}
                    sortField={lazyState.sortField}
                    sortOrder={lazyState.sortOrder}
                >
                    <Column header="Actions" body={actionButtonBodyTemplate} style={{ width: '130px' }} />
                    <Column field="normalized_name" header="Category Name" sortable body={(data: CategoryDto) => TextBodyTemplate(data.name, "Category name", loading)} />
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

const CategoryPage = () => {
    const [filterText, setFilterText] = useState('');
    const [reloading, setReloading] = useState(false);
    const [showModal, setShowModal] = useState(false);
    const [categoryId, setCategoryId] = useState(0);
    const { isGranted } = useSessionStore();

    const breadcrumbs: BreadcrumbItem[] = [
        new BreadcrumbItem('Categories')
    ];

    const getMenuItemsForItem = (item: CategoryDto): MenuItem[] => {
        return [
            {
                label: 'Edit',
                command: (event) => {
                    setCategoryId(item.id ?? 0);
                    setShowModal(true);
                },
                visible: isGranted('Pages.Categories.Edit')
            },
            {
                label: 'Delete',
                command: (event) => {
                    handleDelete(item);
                },
                visible: isGranted('Pages.Categories.Delete')
            }
        ];
    };

    const actionButtons = () => {
        if (!isGranted("Pages.Categories.Create")) {
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
        setCategoryId(0);
        setShowModal(true);
    }

    const handleDelete = (category: CategoryDto) => {
        SwalMessageService.showConfirmation("Are you sure?", `Category "${category.name}" will be deleted`, () => {
            const categoryService = APIClient.getCategoryService();
            categoryService.deleteCategory(category.id ?? 0).then(() => {
                SwalNotifyService.success('Deleted successfully');
                setReloading(!reloading);
            });
        });
    }

    return (
        <>
            <DefaultPage title="Categories" breadcrumbs={breadcrumbs} actionButtons={actionButtons()}>
                <CategoryAdvancedFilter filterText={filterText} setFilterText={setFilterText} />
                <CategoryTable reloading={reloading} filterText={filterText} getMenuItems={getMenuItemsForItem} />
            </DefaultPage>
            <CreateOrEditCategoryModal categoryId={categoryId} show={showModal} handleClose={() => setShowModal(false)} handleSave={() => setReloading(!reloading)} />
        </>
    )
}

export default CategoryPage