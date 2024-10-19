import { useEffect } from "react";
import { Outlet, useLocation } from "react-router-dom"
import { Header } from "../header";
import { AnimatePresence, motion } from "framer-motion";
import { Footer } from "../footer";
import { Sidebar } from "../sidebar";

const pageVariants = {
    initial: { y: '-10%', opacity: 0 },
    enter: { y: '0%', opacity: 1 }
};

const Layout = () => {
    const location = useLocation();
    
    useEffect(() => {
        const bodyClassNames = 'app-layout header-fixed header-tablet-and-mobile-fixed toolbar-enabled';
        const classNamesToAdd = bodyClassNames.split(' ');

        classNamesToAdd.forEach(className => {
            document.body.classList.add(className);
        });

        return () => {
            classNamesToAdd.forEach(className => {
                document.body.classList.remove(className);
            });
        };
    }, []);

    return (
        <>
            <div className="d-flex flex-column flex-root app-root" id="kt_app_root">
                <div className="app-page flex-column flex-column-fluid" id="kt_app_page">
                    <Header />
                    <div className="app-wrapper flex-column flex-row-fluid" id="kt_app_wrapper">
                        <Sidebar />
                        <div className="app-main flex-column flex-row-fluid" id="kt_app_main">
                            <div className="d-flex flex-column flex-column-fluid">
                                <AnimatePresence>
                                    <motion.div
                                        className="animated-box"
                                        key={location.pathname}
                                        initial="initial"
                                        animate="enter"
                                        variants={pageVariants}
                                        transition={{ duration: 0.5 }}
                                    >
                                        <Outlet />
                                    </motion.div>
                                </AnimatePresence>
                            </div>
                            <Footer />
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}

export default Layout