import { useSessionStore } from "@shared/session";
import HeaderUserMenu from "./HeaderUserMenu";
import { Skeleton } from 'primereact/skeleton';
import { Link } from "react-router-dom";

const HeaderTopbar = () => {
    const { user } = useSessionStore();

    return (
        <div className="d-flex align-items-stretch flex-shrink-0">
            <div className="topbar d-flex align-items-stretch flex-shrink-0 me-2 me-md-0">
                <div className="d-flex align-items-center ms-1" title="GoShop">
                    <Link target="_blank" to="https://goshop.gktan.com" className="btn btn-icon btn-custom btn-color-gray-600 btn-active-color-primary w-35px h-35px w-md-40px h-md-40px">
                        <i className="fas fa-shopping-bag fs-2"></i>
                    </Link>
                </div>
                <div className="d-flex align-items-center ms-1" title="Github Link">
                    <Link target="_blank" to="https://github.com/tguankheng016/golang-ecommerce-microservice" className="btn btn-icon btn-custom btn-color-gray-600 btn-active-color-primary w-35px h-35px w-md-40px h-md-40px">
                        <i className="fab fa-github fs-2"></i>
                    </Link>
                </div>
                {
                    !user &&
                    <Skeleton size="3rem" className="mr-2"></Skeleton>
                }
                {
                    user &&
                    <HeaderUserMenu />
                }
            </div>
        </div>
    )
}

export default HeaderTopbar