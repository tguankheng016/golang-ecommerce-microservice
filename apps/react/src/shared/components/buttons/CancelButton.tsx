
interface Props {
    disabled?: boolean;
    onClick?: () => void;
}

const CancelButton = ({ disabled, onClick }: Props) => {
    return (
        <button disabled={disabled} type="button" className="btn btn-light-primary fw-bold" onClick={onClick}>
            Cancel
        </button>
    )
}

export default CancelButton