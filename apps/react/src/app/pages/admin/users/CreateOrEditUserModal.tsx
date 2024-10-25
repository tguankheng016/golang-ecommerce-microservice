import { zodResolver } from "@hookform/resolvers/zod";
import { BusyButton, CancelButton } from "@shared/components/buttons";
import { CustomMessage, ValidationMessage } from "@shared/components/form-validation";
import { DefaultModalProps } from "@shared/components/modals";
import APIClient from "@shared/service-proxies/api-client";
import { CreateOrEditUserDto, CreateUserDto, EditUserDto, RoleDto } from "@shared/service-proxies/identity-service-proxies";
import { SwalNotifyService } from "@shared/sweetalert2";
import { StringHelper } from "@shared/utils";
import { InputText } from "primereact/inputtext";
import { useEffect, useState } from "react";
import { Modal, Tab, Tabs } from "react-bootstrap";
import { useForm } from "react-hook-form";
import { z } from "zod";

interface UserModalProps extends DefaultModalProps {
    userId?: number;
}

class ExtendedRoleDto extends RoleDto {
    isAssigned = false;
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
    const [roles, setRoles] = useState<ExtendedRoleDto[]>([]);

    const { register, reset, handleSubmit, formState: { errors } } = useForm<FormData>({
        resolver: zodResolver(CreateOrEditUserDtoSchema)
    });

    useEffect(() => {
        const abortController = new AbortController();
        const signal = abortController.signal;

        if (show) {
            const userService = APIClient.getUserService();
            const roleService = APIClient.getRoleService();

            const fetchUserById = userService.getUserById(userId ?? 0, signal);
            const fetchRoles = roleService.getAllRoles("", undefined, undefined, undefined, signal);

            Promise.all([fetchUserById, fetchRoles])
                .then(([userRes, rolesRes]) => {
                    setUser(userRes.user ?? new CreateOrEditUserDto());
                    setIsEdit(userRes.user?.id != undefined && userRes.user.id > 0);
                    reset({ ...userRes.user, confirmPassword: "" });

                    const roles = rolesRes.items;
                    if (roles) {
                        const extendedRoles = roles.map((item: RoleDto) => {
                            const extendedItem = new ExtendedRoleDto();
                            Object.assign(extendedItem, item);
                            extendedItem.isAssigned = item.id ? userRes.user?.roleIds?.includes(item.id) ?? false : false;
                            return extendedItem;
                        });
                        setRoles(extendedRoles);
                    }
                });
        }

        return () => {
            abortController.abort();
        };
    }, [show, userId]);

    const handleRoleCheckboxChange = (roleId: number, checked: boolean) => {
        setRoles(prevState => {
            return prevState.map((role) => {
                if (role.id === roleId) {
                    role.isAssigned = checked;
                }
                return role;
            });
        });
    };

    const submitHandler = (data: FormData) => {
        setSaving(true);

        const userService = APIClient.getUserService();

        if (isEdit) {
            // Update user
            const input = EditUserDto.fromJS(data);
            input.id = user.id;
            input.roleIds = roles.filter(x => x.isAssigned).map(x => x.id ?? 0);

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
            input.roleIds = roles.filter(x => x.isAssigned).map(x => x.id ?? 0);

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
        resetForm();
        handleClose();
    }

    const resetForm = () => {
        setUser(new CreateOrEditUserDto());
        setIsEdit(false);
        setSaving(false);
        setRoles([]);
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
                id="kt_create_or_edit_user_form"
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
                    <Tabs defaultActiveKey="general">
                        <Tab eventKey="general" title="General" className="p-3 pt-6">
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
                        </Tab>
                        <Tab eventKey="roles" title="Roles" className="p-3 pt-6">
                            <div className="row">
                                {roles.map((role, index) => (
                                    <div key={role.name}>
                                        <label className="form-check form-check-custom form-check-solid py-2">
                                            <input
                                                id={`User_${role.name}`}
                                                type="checkbox"
                                                name={role.name}
                                                checked={role.isAssigned}
                                                onChange={(e) => handleRoleCheckboxChange(role.id ?? 0, e.target.checked)}
                                                className="form-check-input"
                                            />
                                            <span className="fw-semibold ps-2 fs-6">
                                                {role.name}
                                            </span>
                                        </label>
                                    </div>
                                ))}
                            </div>
                        </Tab>
                    </Tabs>
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