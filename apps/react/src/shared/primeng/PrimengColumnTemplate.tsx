import React from 'react'

export const TextBodyTemplate = (content: string | undefined, header: string): JSX.Element => {
    return (
        <>
            <span className="p-column-title">{header}</span>
            {content}
        </>
    );
}
