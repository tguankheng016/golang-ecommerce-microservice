import { BreadcrumbItem, DefaultPage } from "@app/components/layout";
import { AdvancedFilter } from "@shared/components/advanced-filter";
import { PrimengTableHelper, TextBodyTemplate, useDataTable } from "@shared/primeng";
import { APIClient } from "@shared/service-proxies";
import { CategoryDto, ProductDto } from "@shared/service-proxies/product-service-proxies";
import { useSessionStore } from "@shared/session";
import { SwalMessageService, SwalNotifyService } from "@shared/sweetalert2";
import { Column, ColumnBodyOptions } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { MenuItem } from "primereact/menuitem";
import { Paginator } from "primereact/paginator";
import { TieredMenu } from "primereact/tieredmenu";
import { useEffect, useRef, useState } from "react";
import CreateOrEditProductModal from "./CreateOrEditProductModal";
import { Skeleton } from "primereact/skeleton";
import { Dropdown } from "primereact/dropdown";

interface ProductTableProps {
    categoryIdFilter: number | undefined;
    filterText: string | undefined;
    reloading: boolean;
    getMenuItems: (item: ProductDto) => MenuItem[];
}

const ProductTable = ({ filterText, categoryIdFilter, reloading, getMenuItems }: ProductTableProps) => {
    const menuRefs = useRef<(TieredMenu | null)[]>([]);
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
    const [products, setProducts] = useState<ProductDto[]>([]);

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        loadLazyData(signal);

        return () => {
            abortController.abort();
        };
    }, [lazyState, filterText, categoryIdFilter, reloading]);

    const loadLazyData = (signal: AbortSignal) => {
        const productService = APIClient.getProductService();

        const loadingTimer = setTimeout(() => {
            setLoading(true);
        }, 200);

        productService.getProducts(
            lazyState.rows,
            lazyState.first,
            PrimengTableHelper.getSorting(lazyState),
            filterText,
            categoryIdFilter ?? 0,
            signal
        ).then((res) => {
            setProducts(res.items ?? []);
            setTotalRecords(res.totalCount ?? 0);
        }).finally(() => {
            clearTimeout(loadingTimer);
            setLoading(false);
        });
    };

    const handleButtonClick = (event: React.MouseEvent<HTMLButtonElement>, record: ProductDto, index: number) => {
        setMenuItems(getMenuItems(record));

        if (menuRefs.current[index]) {
            menuRefs.current[index].toggle(event);
        }
    };

    const actionButtonBodyTemplate = (rowData: ProductDto, options: ColumnBodyOptions) => {
        if (loading) {
            return <Skeleton></Skeleton>;
        }

        const assignMenusRef = (ref: TieredMenu | null) => {
            menuRefs.current[options.rowIndex] = ref;
        }

        return (
            <div className="btn-group dropdown">
                <button className="dropdown-toggle btn btn-sm btn-primary" onClick={(e) => handleButtonClick(e, rowData, options.rowIndex)}>
                    <i className="fa fa-cog"></i>
                    <span className="caret"></span>
                    Actions
                </button>
                <TieredMenu model={menuItems} popup ref={assignMenusRef} appendTo={document.body} />
            </div>
        );
    }

    return (
        <div className="row align-items-center mx-0">
            <div className="primeng-datatable-container">
                <DataTable
                    value={products}
                    lazy
                    rows={lazyState.rows}
                    onSort={onSort}
                    sortField={lazyState.sortField}
                    sortOrder={lazyState.sortOrder}
                >
                    <Column header="Actions" body={actionButtonBodyTemplate} style={{ width: '130px' }} />
                    <Column field="p.normalized_name" header="Product Name" sortable body={(data: ProductDto) => TextBodyTemplate(data.name, "Product Name", loading)} />
                    <Column field="p.normalized_description" header="Product Description" sortable body={(data: ProductDto) => TextBodyTemplate(data.description, "Product Descripton", loading)} />
                    <Column field="c.normalized_name" header="Category" sortable body={(data: ProductDto) => TextBodyTemplate(data.categoryName, "Category", loading)} />
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

const ProductPage = () => {
    const [filterText, setFilterText] = useState('');
    const [reloading, setReloading] = useState(false);
    const [showModal, setShowModal] = useState(false);
    const [productId, setProductId] = useState(0);
    const [categoryIdFilter, setCategoryIdFilter] = useState<(number | undefined)>(undefined);
    const [categories, setCategories] = useState<CategoryDto[]>([]);
    const { isGranted } = useSessionStore();

    const breadcrumbs: BreadcrumbItem[] = [
        new BreadcrumbItem('Products')
    ];

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        const categoryService = APIClient.getCategoryService();

        categoryService.getCategories(undefined, undefined, undefined, undefined, signal)
            .then((res) => {
                setCategories(res.items ?? []);
            });

        return () => {
            abortController.abort();
        };
    }, []);

    const getMenuItemsForItem = (item: ProductDto): MenuItem[] => {
        return [
            {
                label: 'Edit',
                command: (event) => {
                    setProductId(item.id ?? 0);
                    setShowModal(true);
                },
                visible: isGranted('Pages.Products.Edit')
            },
            {
                label: 'Delete',
                command: (event) => {
                    handleDelete(item);
                },
                visible: isGranted('Pages.Products.Delete')
            }
        ];
    };

    const actionButtons = () => {
        if (!isGranted("Pages.Products.Create")) {
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
        setProductId(0);
        setShowModal(true);
    }

    const handleDelete = (product: ProductDto) => {
        SwalMessageService.showConfirmation("Are you sure?", `Product "${product.name}" will be deleted`, () => {
            const productService = APIClient.getProductService();
            productService.deleteProduct(product.id ?? 0).then(() => {
                SwalNotifyService.success('Deleted successfully');
                setReloading(!reloading);
            });
        });
    }

    const handleResetFilters = () => {
        setFilterText("");
        setCategoryIdFilter(undefined);
    }

    return (
        <>
            <DefaultPage title="Products" breadcrumbs={breadcrumbs} actionButtons={actionButtons()}>
                <AdvancedFilter filterText={filterText} setFilterText={setFilterText} onResetFilters={handleResetFilters}>
                    <div className="col-md-6 mb-5">
                        <label className="form-label" htmlFor="CategoryIdFilter">
                            Category
                        </label>
                        <Dropdown
                            id="CategoryIdFilter"
                            name="CategoryIdFilter"
                            value={categoryIdFilter}
                            onChange={(e) => setCategoryIdFilter(e.value)}
                            options={categories}
                            optionLabel="name"
                            optionValue="id"
                            showClear={true}
                        />
                    </div>
                </AdvancedFilter>
                <ProductTable reloading={reloading} filterText={filterText} categoryIdFilter={categoryIdFilter} getMenuItems={getMenuItemsForItem} />
            </DefaultPage>
            <CreateOrEditProductModal productId={productId} show={showModal} handleClose={() => setShowModal(false)} handleSave={() => setReloading(!reloading)} />
        </>
    )
}

export default ProductPage