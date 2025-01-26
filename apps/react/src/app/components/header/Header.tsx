import HeaderLogo from "./HeaderLogo"
import HeaderMobileToggle from "./HeaderMobileToggle"
import HeaderTopbar from "./HeaderTopbar"
import HeaderToggle from "./HeaderToggle"

const Header = () => {
    return (
        <div id="kt_app_header" className="app-header d-flex flex-column flex-stack">
            <div className="d-flex align-items-center flex-stack flex-grow-1">
                <div className="app-header-logo d-flex align-items-center flex-stack px-lg-11 mb-2" id="kt_app_header_logo">
                    <HeaderMobileToggle />
                    <HeaderLogo />
                    <HeaderToggle />
                </div>
                <div className="app-navbar flex-grow-1 justify-content-end me-md-12" id="kt_app_header_navbar">
                    <HeaderTopbar />
                </div>
            </div>
        </div>
    )
}

export default Header