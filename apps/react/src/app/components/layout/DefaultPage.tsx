import { Fragment, ReactNode } from "react";
import { Link } from "react-router-dom";

export class BreadcrumbItem {
    text: string;
    routerLink?: string;

    constructor(text: string, routerLink?: string) {
        this.text = text;
        this.routerLink = routerLink;
    }

    isLink(): boolean {
        return !!this.routerLink;
    }
}

interface Props {
    title: string;
    breadcrumbs: BreadcrumbItem[];
    children: ReactNode;
    actionButtons?: ReactNode;
}

const DefaultPage = ({ title, breadcrumbs, children, actionButtons }: Props) => {
    return (
        <>
            <div id="kt_app_toolbar" className="app-toolbar pt-8">
                <div id="kt_app_toolbar_container" className="app-container container-fluid d-flex align-items-stretch">
                    <div className="app-toolbar-wrapper d-flex flex-stack flex-wrap gap-4 w-100">
                        <div className="page-title d-flex flex-column gap-1 me-3 mb-2">
                            <ul className="breadcrumb breadcrumb-separatorless fw-semibold mb-6">
                                <li className="breadcrumb-item text-white opacity-75">
                                    <Link to="/" className="text-white text-hover-primary">
                                        <i className="ki-duotone ki-home fs-3 text-gray-400 me-n1"></i>
                                    </Link>
                                </li>
                                {
                                    breadcrumbs && breadcrumbs.length > 0 &&
                                    <li className="breadcrumb-item">
                                        <i className="ki-duotone ki-right fs-4 text-gray-700 mx-n1"></i>
                                    </li>
                                }
                                {
                                    breadcrumbs.map((item, index) => (
                                        <Fragment key={'breadcrum-item' + index}>
                                            <li className="breadcrumb-item text-gray-700 fw-bold lh-1">
                                                {
                                                    item.isLink() &&
                                                    <a
                                                        className="text-muted cursor-pointer"
                                                    >
                                                        {item.text}
                                                    </a>
                                                }
                                                {
                                                    item.isLink() ||
                                                    <span>{item.text}</span>
                                                }
                                            </li>
                                            {
                                                index !== breadcrumbs.length - 1 &&
                                                <li className="breadcrumb-item">
                                                    <i className="ki-duotone ki-right fs-4 text-gray-700 mx-n1"></i>
                                                </li>
                                            }
                                        </Fragment>
                                    ))
                                }
                            </ul>
                            <h1 className="page-heading d-flex flex-column justify-content-center text-dark fw-bolder fs-1 lh-0">
                                {title}
                            </h1>
                        </div>
                        {actionButtons}
                    </div>
                </div>
            </div>
            <div id="kt_app_content" className="app-content d-flex flex-column flex-root flex-column-fluid">
                <div id="kt_app_content_container" className="app-container container-fluid d-flex flex-column flex-root">
                    {children}
                </div>
            </div>
        </>
    )
}

export default DefaultPage