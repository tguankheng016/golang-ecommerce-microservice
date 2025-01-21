import { useSessionStore } from "@shared/session";
import { Link } from "react-router-dom";
import HeaderUserMenu from "./HeaderUserMenu";
import HeaderCart from "./HeaderCart";

const HeaderTopbar = () => {
    const { user } = useSessionStore();

    return (
        <div className="d-flex align-items-stretch flex-shrink-0">
            <div className="topbar d-flex align-items-stretch flex-shrink-0">
                {
                    (!user || !user.id) &&
                    <Link to='/account/login' className="menu-item menu-lg-down-accordion me-lg-1">
                        <span className="menu-link py-3">
                            <span className="menu-title fs-5 fw-bold">Login</span>
                        </span>
                    </Link>
                }
                {
                    (user && user.id) &&
                    <>
                        <HeaderCart />
                        <HeaderUserMenu />
                    </>
                }
            </div>
        </div>
    )
}

export default HeaderTopbar