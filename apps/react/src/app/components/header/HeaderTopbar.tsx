import useSessionStore from "@shared/session/session-store";
import HeaderUserMenu from "./HeaderUserMenu";
import { Skeleton } from 'primereact/skeleton';

const HeaderTopbar = () => {
    const { user } = useSessionStore();
    
    return (
        <div className="d-flex align-items-stretch flex-shrink-0">
            <div className="topbar d-flex align-items-stretch flex-shrink-0 me-2 me-md-0">
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