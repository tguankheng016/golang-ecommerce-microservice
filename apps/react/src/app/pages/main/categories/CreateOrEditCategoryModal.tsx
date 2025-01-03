import { zodResolver } from "@hookform/resolvers/zod";
import { BusyButton, CancelButton } from "@shared/components/buttons";
import { CustomMessage, ValidationMessage } from "@shared/components/form-validation";
import { DefaultModalProps } from "@shared/components/modals";
import APIClient from "@shared/service-proxies/api-client";
import { CreateOrEditCategoryDto, HumaCreateCategoryRequestBody, HumaUpdateCategoryRequestBody } from "@shared/service-proxies/product-service-proxies";
import { SwalNotifyService } from "@shared/sweetalert2";
import { InputText } from "primereact/inputtext";
import { useEffect, useState } from "react";
import { Modal } from "react-bootstrap";
import { useForm } from "react-hook-form";
import { z } from 'zod';

interface CategoryModalProps extends DefaultModalProps {
    categoryId?: number;
}

const CreateOrEditCategoryDtoSchema = z.object({
    name: z.string().min(1, { message: CustomMessage.required }),
});

type FormData = z.infer<typeof CreateOrEditCategoryDtoSchema>;

const CreateOrEditCategoryModal = ({ categoryId, show, handleClose, handleSave }: CategoryModalProps) => {
    const [saving, setSaving] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const [loading, setLoading] = useState(false);
    const [category, setCategory] = useState<CreateOrEditCategoryDto>(new CreateOrEditCategoryDto());
    const { register, reset, handleSubmit, formState: { errors } } = useForm<FormData>({
        resolver: zodResolver(CreateOrEditCategoryDtoSchema)
    });

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        if (show) {
            setLoading(true);
            const categoryService = APIClient.getCategoryService();

            categoryService.getCategoryById(categoryId ?? 0, signal)
                .then((res) => {
                    setCategory(res.category ?? new CreateOrEditCategoryDto());
                    setIsEdit(res.category?.id != undefined && res.category.id > 0);
                    reset(res.category);
                }).finally(() => {
                    setLoading(false);
                });
        }

        return () => {
            abortController.abort();
        };
    }, [show, categoryId]);

    const submitHandler = (data: FormData) => {
        setSaving(true);

        const categoryService = APIClient.getCategoryService();

        if (isEdit) {
            // Update category
            const input = HumaUpdateCategoryRequestBody.fromJS(data);
            input.id = category.id;

            categoryService.updateCategory(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        } else {
            // Create new category
            const input = HumaCreateCategoryRequestBody.fromJS(data);
            input.id = category.id;

            categoryService.createCategory(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        }
    };

    const resetForm = () => {
        setCategory(new CreateOrEditCategoryDto());
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
                id="kt_create_or_edit_category_form"
            >
                <div className="modal-header">
                    <h5 className="modal-title">
                        {isEdit ? (
                            <span>Edit Category: {category.name}</span>
                        ) : (
                            <span>Create New Category</span>
                        )}
                    </h5>
                    <button type="button" className="btn-close" onClick={closeHandler} aria-label="Close"></button>
                </div>
                <div className="modal-body">
                    <div className="mb-5">
                        <label className="form-label required">Category name</label>
                        <InputText
                            {...register('name')}
                            type="text"
                        />
                        <ValidationMessage errorMessage={errors?.name?.message} />
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

export default CreateOrEditCategoryModal