import { zodResolver } from "@hookform/resolvers/zod";
import { BusyButton, CancelButton } from "@shared/components/buttons";
import { CustomMessage, ValidationMessage } from "@shared/components/form-validation";
import { DefaultModalProps } from "@shared/components/modals";
import APIClient from "@shared/service-proxies/api-client";
import { CategoryDto, CreateOrEditProductDto, HumaCreateProductRequestBody, HumaUpdateProductRequestBody } from "@shared/service-proxies/product-service-proxies";
import { SwalNotifyService } from "@shared/sweetalert2";
import { Dropdown } from "primereact/dropdown";
import { InputNumber } from "primereact/inputnumber";
import { InputText } from "primereact/inputtext";
import { useEffect, useState } from "react";
import { Modal } from "react-bootstrap";
import { useForm } from "react-hook-form";
import { z } from 'zod';

interface ProductModalProps extends DefaultModalProps {
    productId?: number;
}

const CreateOrEditProductDtoSchema = z.object({
    name: z.string().min(1, { message: CustomMessage.required }),
    description: z.string().min(1, { message: CustomMessage.required }),
    price: z.string().min(1, { message: CustomMessage.required }),
    stockQuantity: z.coerce.number().int().min(0, { message: CustomMessage.required }),
    categoryId: z.coerce.number().int().min(1, { message: CustomMessage.required }),
});

type FormData = z.infer<typeof CreateOrEditProductDtoSchema>;

const CreateOrEditProductModal = ({ productId, show, handleClose, handleSave }: ProductModalProps) => {
    const [saving, setSaving] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const [loading, setLoading] = useState(false);
    const [product, setProduct] = useState<CreateOrEditProductDto>(new CreateOrEditProductDto());
    const [categories, setCategories] = useState<CategoryDto[]>([]);
    const { register, reset, handleSubmit, setValue, formState: { errors } } = useForm<FormData>({
        resolver: zodResolver(CreateOrEditProductDtoSchema)
    });

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        if (show) {
            setLoading(true);
            const productService = APIClient.getProductService();
            const categoryService = APIClient.getCategoryService();

            const categories$ = categoryService.getCategories(undefined, undefined, undefined, undefined, signal);
            const product$ = productService.getProductById(productId ?? 0, signal);

            Promise.all([categories$, product$])
                .then(res => {
                    setCategories(res[0].items ?? []);
                    setProduct(res[1].product ?? new CreateOrEditProductDto());
                    setIsEdit(res[1].product?.id != undefined && res[1].product.id > 0);
                    reset(res[1].product);
                });
        }

        return () => {
            abortController.abort();
        };
    }, [show, productId]);

    const submitHandler = (data: FormData) => {
        setSaving(true);

        const productService = APIClient.getProductService();

        if (isEdit) {
            // Update product
            const input = HumaUpdateProductRequestBody.fromJS(data);
            input.id = product.id;

            productService.updateProduct(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        } else {
            // Create new product
            const input = HumaCreateProductRequestBody.fromJS(data);
            input.id = 0;

            productService.createProduct(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        }
    };

    const resetForm = () => {
        setProduct(new CreateOrEditProductDto());
        setIsEdit(false);
        setSaving(false);
    }

    const closeHandler = () => {
        resetForm();
        handleClose();
    }

    return (
        <Modal
            show={show}
            onHide={handleClose}
            backdrop="static"
            keyboard={false}
            size="lg"
        >
            <form
                onSubmit={handleSubmit(data => submitHandler(data))}
                className="form w-100 fv-plugins-bootstrap5 fv-plugins-framework"
                id="kt_create_or_edit_product_form"
            >
                <div className="modal-header">
                    <h5 className="modal-title">
                        {isEdit ? (
                            <span>Edit Product: {product.name}</span>
                        ) : (
                            <span>Create New Product</span>
                        )}
                    </h5>
                    <button type="button" className="btn-close" onClick={closeHandler} aria-label="Close"></button>
                </div>
                <div className="modal-body">
                    <div className="mb-5">
                        <label className="form-label required">Product name</label>
                        <InputText
                            {...register('name')}
                            type="text"
                        />
                        <ValidationMessage errorMessage={errors?.name?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Product Description</label>
                        <InputText
                            {...register('description')}
                            type="text"
                        />
                        <ValidationMessage errorMessage={errors?.description?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Price</label>
                        <InputNumber
                            value={parseFloat(product.price)}
                            onValueChange={(e) => {
                                setValue('price', e.value?.toString() ?? "");
                                setProduct((prevProduct) => (new CreateOrEditProductDto({ ...prevProduct, price: e.value?.toString() ?? "" })));
                            }}
                            useGrouping={true}
                            minFractionDigits={2}
                            maxFractionDigits={2}
                        />
                        <ValidationMessage errorMessage={errors?.price?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Stock Quantity</label>
                        <InputNumber
                            value={product.stockQuantity}
                            onValueChange={(e) => {
                                setValue('stockQuantity', e.value ?? 0);
                                setProduct((prevProduct) => (new CreateOrEditProductDto({ ...prevProduct, stockQuantity: e.value ?? 0 })));
                            }}
                            useGrouping={true}
                            minFractionDigits={0}
                            maxFractionDigits={0}
                        />
                        <ValidationMessage errorMessage={errors?.stockQuantity?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Category</label>
                        <Dropdown
                            value={product.categoryId}
                            options={categories}
                            optionLabel="name"
                            optionValue="id"
                            onChange={(e) => {
                                setValue('categoryId', e.value);
                                setProduct((prevProduct) => (new CreateOrEditProductDto({ ...prevProduct, categoryId: e.value })));
                            }}
                        />
                        <ValidationMessage errorMessage={errors?.categoryId?.message} />
                    </div>
                </div>
                <div className="modal-footer">
                    <CancelButton disabled={saving} onClick={closeHandler} />
                    <BusyButton
                        className="btn btn-primary fw-bold"
                        busyIf={saving}
                    >
                        <i className="fa fa-save me-2"></i>
                        <span>Save</span>
                    </BusyButton>
                </div>
            </form>
        </Modal>
    )
}

export default CreateOrEditProductModal