const Footer = () => {
    return (
        <div id="kt_app_footer" className="app-footer align-items-center justify-content-center justify-content-md-between flex-column flex-md-row py-3">
            <div className="text-dark order-2 order-md-1">
                <span className="text-muted fw-semibold me-1">2023&copy;</span>
                <a href="https://keenthemes.com" target="_blank" className="text-gray-800 text-hover-primary">Keenthemes</a>
            </div>
            <ul className="menu menu-gray-600 menu-hover-primary fw-semibold order-1">
                <li className="menu-item">
                    <a href="https://keenthemes.com" target="_blank" className="menu-link px-2">About</a>
                </li>
                <li className="menu-item">
                    <a href="https://devs.keenthemes.com" target="_blank" className="menu-link px-2">Support</a>
                </li>
                <li className="menu-item">
                    <a href="https://keenthemes.com/products/saul-html-pro" target="_blank" className="menu-link px-2">Purchase</a>
                </li>
            </ul>
        </div>
    )
}

export default Footer