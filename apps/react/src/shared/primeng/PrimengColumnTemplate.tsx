import moment from 'moment';

export const TextBodyTemplate = (content: string | undefined, header: string): JSX.Element => {
    return (
        <>
            <span className="p-column-title">{header}</span>
            {content}
        </>
    );
}

export const DateTimeBodyTemplate = (content: moment.Moment | undefined, header: string, format: string = "YYYY-MM-DD HH:mm:ss"): JSX.Element => {
    const formattedDate = content ? moment(content).format(format) : '';

    return (
        <>
            <span className="p-column-title">{header}</span>
            {formattedDate}
        </>
    );
}
