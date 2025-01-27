import DefaultPage, { BreadcrumbItem } from "@app/components/layout/DefaultPage"

const HomePage = () => {
    const breadcrumbs: BreadcrumbItem[] = [
        new BreadcrumbItem('About')
    ];

    return (
        <DefaultPage title="About" breadcrumbs={breadcrumbs}>
            <p>
                This is a practical and imaginary eCommerce system built
                with <strong>Go</strong> as the backend and <strong>React</strong> for the
                frontend. The system is designed using <strong>Microservices Architecture</strong>,
                <strong>Vertical Slice Architecture</strong>, and <strong>Clean Architecture</strong> principles.
                The goal is to demonstrate a scalable, maintainable, and testable approach for
                building modern eCommerce applications.
            </p>
            <h3>Features</h3>
            <ul>
                <li>
                    <span>
                        Using <b>Nats Jetstream</b> on top of <b>Watermill</b> for
                        Asynchronous Messaging
                    </span>
                </li>
                <li>
                    <span>
                        Using <b>gRPC</b> for <b>internal communications</b> between
                        services
                    </span>
                </li>
                <li>
                    <span> Using <b>Uber Zap</b> for structured logging </span>
                </li>
                <li>
                    <span> Using <b>Uber Fx</b> for dependency injection </span>
                </li>
                <li>
                    <span> Using <b>Huma</b> and <b>Chi</b> to handle requests </span>
                </li>
                <li>
                    <span>
                        Using <b>Go Validating</b> to validate input requests
                    </span>
                </li>
                <li>
                    <span> Using <b>Goose</b> to handle migrations </span>
                </li>
                <li>
                    <span> Using <b>Viper</b> to manage configuration </span>
                </li>
                <li>
                    <span>
                        Using <b>PostgreSQL</b> and <b>MongoDB</b> as databases
                    </span>
                </li>
                <li>
                    <span> Using <b>Stoplight Elements</b> for OpenAPI documentation </span>
                </li>
                <li>
                    <span>
                        Using <b>Nswag</b> for generating Typescript client code
                    </span>
                </li>
                <li>
                    <span> Using <b>Redis</b> for distributed caching </span>
                </li>
                <li>
                    <span>
                        Using <b>OpenTelemetry</b> and <b>Jaeger</b> for distributed
                        tracing.
                    </span>
                </li>
                <li>
                    <span>
                        Using <b>Testify</b>, <b>TestContainers</b> for unit and
                        integration testing
                    </span>
                </li>
            </ul>
            <h3>Other Relevant Sites</h3>
            <ul>
                <li>
                    <span>GoShop:</span>&nbsp;
                    <a href="https://goshop.gktan.com" target="_blank"
                    >https://goshop.gktan.com</a
                    >
                </li>
                <li>
                    <span>OAuth:</span>&nbsp;
                    <a href="https://auth.gktan.com" target="_blank"
                    >https://auth.gktan.com</a
                    >
                </li>
            </ul>
            <h3>Source code</h3>
            <p>
                This project is developed open source on Github.&nbsp;
                <a
                    href="https://github.com/tguankheng016/golang-ecommerce-microservice"
                    target="_blank"
                >https://github.com/tguankheng016/golang-ecommerce-microservice</a
                >
            </p>
        </DefaultPage>
    )
}

export default HomePage