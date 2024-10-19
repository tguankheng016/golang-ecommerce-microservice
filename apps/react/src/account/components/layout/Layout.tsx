import { Outlet } from "react-router-dom"
import logo from '@assets/media/mail.svg';

const Layout = () => {
    return (
        <div className="d-flex flex-column flex-root" id="kt_app_root">
            <div className="d-flex flex-column flex-lg-row flex-column-fluid">
                <div className="d-flex flex-column flex-lg-row-auto bg-primary w-xl-600px positon-xl-relative">
                    <div className="d-flex flex-column position-xl-fixed top-0 bottom-0 w-xl-600px scroll-y">
                        <div className="d-flex flex-row-fluid flex-column text-center p-5 p-lg-10 pt-lg-20">
                            <a href="./" className="py-2 py-lg-20">
                                <img alt="Logo" src={logo} className="h-40px h-lg-50px" />
                            </a>
                            <h1 className="d-none d-lg-block fw-bold text-white fs-2qx pb-5 pb-md-10">Welcome to Saul HTML Free</h1>
                            <p className="d-none d-lg-block fw-semibold fs-2 text-white">Plan your blog post by choosing a topic creating
                            <br />an outline and checking facts</p>
                        </div>
                        <div className="auth-page-bg d-none d-lg-block d-flex flex-row-auto bgi-no-repeat bgi-position-x-center bgi-size-contain bgi-position-y-bottom min-h-100px min-h-lg-350px"></div>
                    </div>
                </div>
                <div className="d-flex flex-column flex-lg-row-fluid py-10">
                    <div className="d-flex flex-center flex-column flex-column-fluid">
                        <div className="w-lg-500px p-10 p-lg-15 mx-auto w-100 w-md-unset">
                            <Outlet />
                        </div>
                    </div>
                    <div className="d-flex flex-center flex-wrap fs-6 p-5 pb-0">
                        <div className="d-flex flex-center fw-semibold fs-6">
                            <a href="https://keenthemes.com" className="text-muted text-hover-primary px-2" target="_blank">About</a>
                            <a href="https://devs.keenthemes.com" className="text-muted text-hover-primary px-2" target="_blank">Support</a>
                            <a href="https://keenthemes.com/products/saul-html-pro" className="text-muted text-hover-primary px-2" target="_blank">Purchase</a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Layout