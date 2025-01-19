import { zodResolver } from "@hookform/resolvers/zod";
import { BusyButton, CancelButton } from "@shared/components/buttons";
import { CustomMessage, ValidationMessage } from "@shared/components/form-validation";
import { DefaultModalProps } from "@shared/components/modals";
import { APIClient } from "@shared/service-proxies";
import { CategoryDto, HumaCreateProductRequestBody, HumaUpdateProductRequestBody } from "@shared/service-proxies/product-service-proxies";
import { SwalNotifyService } from "@shared/sweetalert2";
import { Dropdown } from "primereact/dropdown";
import { InputNumber } from "primereact/inputnumber";
import { InputText } from "primereact/inputtext";
import { classNames } from "primereact/utils";
import { useEffect, useRef, useState } from "react";
import { Modal } from "react-bootstrap";
import { Controller, useForm } from "react-hook-form";
import { z } from 'zod';

interface ProductModalProps extends DefaultModalProps {
    productId?: number;
}

const CreateOrEditProductDtoSchema = z.object({
    id: z.number().nullable().optional(),
    name: z.string().min(1, { message: CustomMessage.required }),
    description: z.string().min(1, { message: CustomMessage.required }),
    price: z.number({ invalid_type_error: CustomMessage.required }).min(0),
    stockQuantity: z.number({ invalid_type_error: CustomMessage.required }).min(0),
    categoryId: z.coerce.number().int().min(1, { message: CustomMessage.required }),
});

type FormData = z.infer<typeof CreateOrEditProductDtoSchema>;

const CreateOrEditProductModal = ({ productId, show, handleClose, handleSave }: ProductModalProps) => {
    const [saving, setSaving] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const [loading, setLoading] = useState(false);
    const [categories, setCategories] = useState<CategoryDto[]>([]);

    const { control, reset, handleSubmit, getValues, formState: { errors } } = useForm<FormData>({
        resolver: zodResolver(CreateOrEditProductDtoSchema),
        mode: "onTouched",
        defaultValues: {
            name: "",
            description: "",
            price: undefined,
            stockQuantity: undefined,
            categoryId: undefined
        },
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
                    setIsEdit(res[1].product?.id != undefined && res[1].product.id > 0);
                    reset({ ...res[1].product, price: parseFloat(res[1].product.price) });
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
            input.price = input.price.toString();

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
            input.price = input.price.toString();

            productService.createProduct(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        }
    };

    const closeHandler = () => {
        resetForm();
        handleClose();
    }

    const resetForm = () => {
        reset();
        setIsEdit(false);
        setSaving(false);
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
                            <span>Edit Product: {getValues("name")}</span>
                        ) : (
                            <span>Create New Product</span>
                        )}
                    </h5>
                    <button type="button" className="btn-close" onClick={closeHandler} aria-label="Close"></button>
                </div>
                <div className="modal-body">
                    <div className="mb-5">
                        <label className="form-label required">Product name</label>
                        <Controller
                            name="name"
                            control={control}
                            render={({ field, fieldState }) => (
                                <InputText
                                    id={field.name}
                                    {...field}
                                    className={classNames({ 'p-invalid': fieldState.invalid })}
                                    type="text"
                                    autoFocus
                                />
                            )}
                        />
                        <ValidationMessage errorMessage={errors?.name?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Product Description</label>
                        <Controller
                            name="description"
                            control={control}
                            render={({ field, fieldState }) => (
                                <InputText
                                    id={field.name}
                                    {...field}
                                    className={classNames({ 'p-invalid': fieldState.invalid })}
                                    type="text"
                                />
                            )}
                        />
                        <ValidationMessage errorMessage={errors?.description?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Price</label>
                        <Controller
                            name="price"
                            control={control}
                            render={({ field, fieldState }) => (
                                <InputNumber
                                    id={field.name}
                                    value={field.value}
                                    className={classNames({ 'p-invalid': fieldState.invalid })}
                                    onValueChange={(e) => field.onChange(e.value)}
                                    useGrouping={true}
                                    minFractionDigits={2}
                                    maxFractionDigits={2}
                                />
                            )}
                        />
                        <ValidationMessage errorMessage={errors?.price?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Stock Quantity</label>
                        <Controller
                            name="stockQuantity"
                            control={control}
                            render={({ field, fieldState }) => (
                                <InputNumber
                                    id={field.name}
                                    value={field.value}
                                    className={classNames({ 'p-invalid': fieldState.invalid })}
                                    onValueChange={(e) => field.onChange(e.value)}
                                    useGrouping={true}
                                    minFractionDigits={0}
                                    maxFractionDigits={0}
                                />
                            )}
                        />
                        <ValidationMessage errorMessage={errors?.stockQuantity?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Category</label>
                        <Controller
                            name="categoryId"
                            control={control}
                            render={({ field, fieldState }) => (
                                <Dropdown
                                    id={field.name}
                                    value={field.value}
                                    className={classNames({ 'p-invalid': fieldState.invalid })}
                                    onChange={(e) => field.onChange(e.value)}
                                    options={categories}
                                    optionLabel="name"
                                    optionValue="id"
                                />
                            )}
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