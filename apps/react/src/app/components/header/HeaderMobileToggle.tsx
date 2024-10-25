import { useLayoutStore, useThemeStore } from "@shared/theme";

const HeaderMobileToggle = () => {
    const { isDarkMode } = useThemeStore();
    const { isExpanded, setExpanded } = useLayoutStore();

    const handleToggle = () => {
        console.log('toggled');
        setExpanded();
    };

    return (
        <div
            id="kt_app_sidebar_mobile_toggle"
            className={`btn btn-icon btn-active-color-primary w-35px h-35px ms-3 me-2 d-flex d-lg-none${isExpanded ? ' active' : ''}`}
            onClick={handleToggle}
        >
            <i className="ki-duotone ki-abstract-14 fs-2">
                <span className="path1"></span>
                <span className="path2"></span>
            </i>
        </div>
    )
}

export default HeaderMobileToggle