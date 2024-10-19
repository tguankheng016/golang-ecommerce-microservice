import welcomeImg from '@assets/media/saul-welcome.png';
import SidebarMenu from './SidebarMenu';
import useLayoutStore from '@shared/theme/layout-store';
import { useEffect, useRef } from 'react';

const Sidebar = () => {
    const sideBarRef = useRef<HTMLDivElement>(null);
    const { isExpanded, isMobileView, setExpanded } = useLayoutStore();

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (isMobileView && isExpanded && sideBarRef.current && !sideBarRef.current.contains(event.target as Node)) {
                setExpanded();
            }
        };

        if (isMobileView) {
            // Attach the event listener
            document.addEventListener('mousedown', handleClickOutside);

            if (sideBarRef.current) {
                sideBarRef.current.classList.add('drawer', 'drawer-start');
                sideBarRef.current.style.setProperty('width', '250px', 'important');
            }
        } else {
            document.removeEventListener('mousedown', handleClickOutside);
            if (sideBarRef.current) {
                sideBarRef.current.classList.remove('drawer', 'drawer-start');
                sideBarRef.current.style.removeProperty('width');
            }
            removeOverlay();
        }

        // Cleanup the event listener on component unmount
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isMobileView, isExpanded]);

    useEffect(() => {
        if (!isMobileView) {
            return;
        }

        if (isExpanded) {
            document.body.setAttribute('data-kt-drawer', 'on');
            document.body.setAttribute('data-kt-drawer-app-sidebar', 'on');
            sideBarRef.current?.classList.add('drawer-on');
            const overlayDiv = document.createElement('div');
            overlayDiv.id = 'sidebarOverlayDiv';
            overlayDiv.classList.add('drawer-overlay');
            overlayDiv.style.zIndex = '105';
            document.body.appendChild(overlayDiv);
        } else {
            document.body.removeAttribute('data-kt-drawer');
            document.body.removeAttribute('data-kt-drawer-app-sidebar');
            sideBarRef.current?.classList.remove('drawer-on');
            removeOverlay();
        }
    }, [isExpanded]);

    useEffect(() => {
        // Hide the sidebar when the location changes in mobile view
        if (isExpanded && isMobileView) {
            setExpanded();
        }
    }, [location]);

    const removeOverlay = () => {
        const overlayDiv = document.getElementById('sidebarOverlayDiv');
        if (overlayDiv) {
            document.body.removeChild(overlayDiv);
        }
    };
    
    return (
        <div 
            id="kt_app_sidebar" 
            className="app-sidebar flex-column"
            ref={sideBarRef}
        >
            <div 
                className="d-flex flex-column justify-content-between h-100 hover-scroll-overlay-y my-2 d-flex flex-column" 
                id="kt_app_sidebar_main"
            >
                <div 
                    id="#kt_app_sidebar_menu" 
                    className="flex-column-fluid menu menu-sub-indention menu-column menu-rounded menu-active-bg mb-7"
                >
                    <SidebarMenu />
                </div>
                <div className="app-sidebar-project-default app-sidebar-project-minimize text-center min-h-lg-400px flex-column-auto d-flex flex-column justify-content-end" id="kt_app_sidebar_footer">
                    <h2 className="fw-bold text-gray-800">Welcome to Saul</h2>
                    <div className="fw-semibold text-gray-700 fs-7 lh-2 px-7 mb-1">Join the movement make a difference.</div>
                    <img className="mx-auto h-150px h-lg-175px mb-4" src={welcomeImg} alt="Welcome Image" />
                    <div className="text-center mb-lg-5 pb-lg-3">
                        <a href="#" className="btn btn-sm btn-dark" data-bs-toggle="modal" data-bs-target="#kt_modal_create_account">Get Started</a>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Sidebar