import React from 'react'

const ChangePermissionWarningBox = () => {
    return (
        <div className="notice d-flex bg-light-warning rounded border-warning border border-dashed mb-9 p-6 mt-5">
            <i className="ki-duotone ki-design-1 fs-2tx text-primary me-4"></i>
            <div className="d-flex flex-stack flex-grow-1">
                <div className="fw-semibold">
                    <div className="fs-6 text-gray-700">
                        If you are changing your own permissions, you may need to refresh page (F5) to take effect of permission changes on your own screen.
                    </div>
                </div>
            </div>
        </div>
    )
}

export default ChangePermissionWarningBox