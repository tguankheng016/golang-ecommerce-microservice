import { zodResolver } from "@hookform/resolvers/zod";
import { BusyButton, CancelButton } from "@shared/components/buttons";
import { CustomMessage, ValidationMessage } from "@shared/components/form-validation";
import APIClient from "@shared/service-proxies/api-client";
import { CreateOrEditUserDto, CreateUserDto, EditUserDto } from "@shared/service-proxies/identity-service-proxies";
import SwalNotifyService from "@shared/sweetalert2/swal-notify";
import StringHelper from "@shared/utils/string-helper";
import { InputText } from "primereact/inputtext";
import { useEffect, useMemo, useState } from "react";
import { Modal } from "react-bootstrap";
import { useForm } from "react-hook-form";
import { z } from "zod";

interface UserModalProps {
    userId?: number;
    show: boolean;
    handleClose: () => void;
    handleSave?: () => void;
}

const CreateOrEditUserDtoSchema = z.object({
    id: z.number().optional(),
    userName: z.string().min(1, { message: CustomMessage.required }),
    firstName: z.string().min(1, { message: CustomMessage.required }),
    lastName: z.string().min(1, { message: CustomMessage.required }),
    email: z.string().email(CustomMessage.invalidEmail),
    password: z.string().optional(),
    confirmPassword: z.string().optional()
}).superRefine(({ id, password, confirmPassword }, ctx) => {
    const isEdit = id !== undefined && id > 0;

    if (isEdit && !StringHelper.notNullOrEmpty(password) && !StringHelper.notNullOrEmpty(confirmPassword)) {
        return ctx;
    }

    if (!isEdit) {
        if (!password || password?.length < 6) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: "Password must be at least 6 characters",
                path: ['password'],
            });
        }

        if (!confirmPassword || confirmPassword?.length < 6) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message: "Password must be at least 6 characters",
                path: ['confirmPassword'],
            });
        }
    }

    if (password !== confirmPassword) {
        ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: "Passwords do not match",
            path: ['confirmPassword'],
        });
    }

    return ctx;
});

type FormData = z.infer<typeof CreateOrEditUserDtoSchema>;

const CreateOrEditUserModal = ({ userId, show, handleClose, handleSave }: UserModalProps) => {
    const [saving, setSaving] = useState(false);
    const [isEdit, setIsEdit] = useState(false);
    const [user, setUser] = useState<CreateOrEditUserDto>(new CreateOrEditUserDto());

    const { register, reset, handleSubmit, formState: { errors } } = useForm<FormData>({
        resolver: zodResolver(CreateOrEditUserDtoSchema)
    });

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        if (show) {
            const userService = APIClient.getUserService();

            userService.getUserById(
                userId ?? 0,
                signal
            ).then((res) => {
                setUser(res.user ?? new CreateOrEditUserDto());
                setIsEdit(res.user?.id != undefined && res.user.id > 0);
                reset({ ...res.user, confirmPassword: "" });
            });
        }

        return () => {
            abortController.abort();
        };
    }, [show, userId]);

    const submitHandler = (data: FormData) => {
        setSaving(true);

        const userService = APIClient.getUserService();

        if (isEdit) {
            // Update user
            const input = EditUserDto.fromJS(data);
            input.id = user.id;
            userService.updateUser(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        } else {
            // Create new user
            const input = CreateUserDto.fromJS(data);
            userService.createUser(input).then((res) => {
                SwalNotifyService.info('Saved successfully');
                closeHandler();
                handleSave?.();
            }).finally(() => {
                setSaving(false);
            });
        }
    };

    const closeHandler = () => {
        setUser(new CreateOrEditUserDto());
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
                id="kt_sign_in_form"
            >
                <div className="modal-header">
                    <h5 className="modal-title">
                        {isEdit ? (
                            <span>Edit User: {user.userName}</span>
                        ) : (
                            <span>Create New User</span>
                        )}
                    </h5>
                    <button type="button" className="btn-close" onClick={closeHandler} aria-label="Close"></button>
                </div>
                <div className="modal-body">
                    <div className="mb-5">
                        <label className="form-label required">First name</label>
                        <InputText
                            {...register('firstName')}
                            type="text"
                        />
                        <ValidationMessage errorMessage={errors?.firstName?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Last name</label>
                        <InputText
                            {...register('lastName')}
                            type="text"
                        />
                        <ValidationMessage errorMessage={errors?.lastName?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Email</label>
                        <InputText
                            {...register('email')}
                            type="text"
                        />
                        <ValidationMessage errorMessage={errors?.email?.message} />
                    </div>
                    <div className="mb-5">
                        <label className="form-label required">Username</label>
                        <InputText
                            {...register('userName')}
                            type="text"
                        />
                        <ValidationMessage errorMessage={errors?.userName?.message} />
                    </div>
                    <div className="mb-5">
                        <label className={`form-label${isEdit ? "" : " required"}`}>Password</label>
                        <InputText
                            {...register('password')}
                            type="password"
                        />
                        <ValidationMessage errorMessage={errors?.password?.message} />
                    </div>
                    <div className="mb-5">
                        <label className={`form-label${isEdit ? "" : " required"}`}>Confirm Password</label>
                        <InputText
                            {...register('confirmPassword')}
                            type="password"
                        />
                        <ValidationMessage errorMessage={errors?.confirmPassword?.message} />
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

export default CreateOrEditUserModal