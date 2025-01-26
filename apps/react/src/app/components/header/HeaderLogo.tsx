import logo from '@assets/media/default-dark.svg';
import logoLight from '@assets/media/default.svg';

const HeaderLogo = () => {
    return (
        <div className="header-logo me-5 me-md-10 flex-grow-1 flex-md-grow-0">
            <a href="/" className="app-sidebar-logo">
                <img alt="Logo" src={logoLight} className="h-30px theme-light-show" />
                <img alt="Logo" src={logo} className="h-30px theme-dark-show" />
            </a>
        </div>
    )
}

export default HeaderLogo