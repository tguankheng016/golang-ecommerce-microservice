import { useLayoutStore } from '@shared/theme';

const HeaderToggle = () => {
    const { isExpanded, setExpanded } = useLayoutStore();

    const handleToggle = () => {
        if (isExpanded) {
            document.body.removeAttribute("data-kt-app-sidebar-minimize");
        } else {
            document.body.setAttribute("data-kt-app-sidebar-minimize", "on");
        }

        setExpanded();
    };

    return (
        <>
            <div
                id="kt_app_sidebar_toggle"
                className={`app-sidebar-toggle btn btn-sm btn-icon btn-color-primary me-n2 d-none d-lg-flex${isExpanded ? ' active' : ''}`}
                data-kt-toggle="true"
                data-kt-toggle-state="active"
                data-kt-toggle-target="body"
                data-kt-toggle-name="app-sidebar-minimize"
                onClick={handleToggle}
            >
                <i className="ki-duotone ki-exit-left fs-2x rotate-180">
                    <span className="path1"></span>
                    <span className="path2"></span>
                </i>
            </div>
        </>
    )
}

export default HeaderToggle