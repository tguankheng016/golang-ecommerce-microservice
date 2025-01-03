import { BreadcrumbItem, DefaultPage } from "@app/components/layout";
import { AdvancedFilter, AdvancedFilterProps } from "@shared/components/advanced-filter";
import { PrimengTableHelper, useDataTable } from "@shared/primeng";
import APIClient from "@shared/service-proxies/api-client";
import { ProductDto } from "@shared/service-proxies/product-service-proxies";
import { useSessionStore } from "@shared/session";
import { SwalMessageService, SwalNotifyService } from "@shared/sweetalert2";
import { Column } from "primereact/column";
import { DataTable } from "primereact/datatable";
import { MenuItem } from "primereact/menuitem";
import { Paginator } from "primereact/paginator";
import { TieredMenu } from "primereact/tieredmenu";
import { useEffect, useRef, useState } from "react";
import CreateOrEditProductModal from "./CreateOrEditProductModal";

interface ProductTableProps {
    filterText: string | undefined;
    reloading: boolean;
    getMenuItems: (item: ProductDto) => MenuItem[];
}

const ProductAdvancedFilter = ({ filterText, setFilterText }: AdvancedFilterProps) => {
    return (
        <AdvancedFilter
            filterText={filterText}
            setFilterText={setFilterText}
        >
        </AdvancedFilter>
    )
}

const ProductTable = ({ filterText, reloading, getMenuItems }: ProductTableProps) => {
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
    const [products, setProducts] = useState<ProductDto[]>([]);

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        loadLazyData(signal);

        return () => {
            abortController.abort();
        };
    }, [lazyState, filterText, reloading]);

    const loadLazyData = (signal: AbortSignal) => {
        const productService = APIClient.getProductService();

        setLoading(true);

        productService.getProducts(
            lazyState.rows,
            lazyState.first,
            PrimengTableHelper.getSorting(lazyState),
            filterText,
            signal
        ).then((res) => {
            setProducts(res.items ?? []);
            setTotalRecords(res.totalCount ?? 0);
            setLoading(false);
        }).finally(() => {
            setLoading(false);
        });
    };

    const handleButtonClick = (event: React.MouseEvent<HTMLButtonElement>, record: ProductDto) => {
        setMenuItems(getMenuItems(record));
        if (menu.current) {
            menu.current.toggle(event);
        }
    };

    const actionButtonBodyTemplate = (rowData: ProductDto) => {
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

    const productNameTemplate = (rowData: ProductDto) => {
        return (
            <>
                <span className="p-column-title">Product name</span>
                <span>{rowData.name}</span>
            </>
        );
    }

    const productDescriptionTemplate = (rowData: ProductDto) => {
        return (
            <>
                <span className="p-column-title">Product description</span>
                <span>{rowData.description}</span>
            </>
        );
    }

    const productCategoryNameTemplate = (rowData: ProductDto) => {
        return (
            <>
                <span className="p-column-title">Category</span>
                <span>{rowData.categoryName}</span>
            </>
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
                    <Column field="p.normalized_name" header="Product Name" sortable body={productNameTemplate} />
                    <Column field="p.normalized_description" header="Product Description" sortable body={productDescriptionTemplate} />
                    <Column field="c.normalized_name" header="Category" sortable body={productCategoryNameTemplate} />
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
    const { isGranted } = useSessionStore();

    const breadcrumbs: BreadcrumbItem[] = [
        new BreadcrumbItem('Products')
    ];

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

    return (
        <>
            <DefaultPage title="Products" breadcrumbs={breadcrumbs} actionButtons={actionButtons()}>
                <ProductAdvancedFilter filterText={filterText} setFilterText={setFilterText} />
                <ProductTable reloading={reloading} filterText={filterText} getMenuItems={getMenuItemsForItem} />
            </DefaultPage>
            <CreateOrEditProductModal productId={productId} show={showModal} handleClose={() => setShowModal(false)} handleSave={() => setReloading(!reloading)} />
        </>
    )
}

export default ProductPage