import { InputText } from 'primereact/inputtext';
import { Dispatch, SetStateAction, useRef, useState } from 'react'
import { useDebounceCallback } from 'usehooks-ts'

export interface AdvancedFilterProps {
    filterText: string;
    setFilterText: Dispatch<SetStateAction<string>>;
    onResetFilters: () => void;
    children?: React.ReactNode;
}

const AdvancedFilter = ({ setFilterText, onResetFilters, children }: AdvancedFilterProps) => {
    const [showAdvancedFilter, setShowAdvancedFilter] = useState(false);
    const setDebouncedFilterText = useDebounceCallback(setFilterText, 300);
    const inputRef = useRef<HTMLInputElement | null>(null);

    const handleResetButtonClick = () => {
        onResetFilters();

        if (inputRef?.current) {
            inputRef.current.value = "";
        }
    }

    return (
        <form className="form">
            <div className="row align-items-center mb-2">
                <div className="col-xl-12">
                    <div className="mb-2 m-form__group align-items-center">
                        <div className="input-group">
                            <div className="input-group-prepend">
                                <button className="btn btn-primary" type="button">
                                    <i className="fas fa-search" aria-label="Search"></i>
                                </button>
                            </div>
                            <InputText
                                ref={inputRef}
                                name="filterText"
                                placeholder="Search ..."
                                onChange={(e) => setDebouncedFilterText(e.target.value)}
                            />
                        </div>
                    </div>
                </div>
            </div>
            {
                showAdvancedFilter &&
                <div className="row">
                    {children}
                    <div className="col-md-12 mt-5 mb-3">
                        <button
                            id="btn-reset-filters"
                            className="btn btn-secondary btn-sm"
                            onClick={handleResetButtonClick}
                            type="button"
                        >
                            Reset
                        </button>
                    </div>
                </div>
            }
            <div className="row mb-2">
                <div className="col-sm-6">
                    {
                        !showAdvancedFilter &&
                        <span
                            className="cursor-pointer text-muted"
                            onClick={() => setShowAdvancedFilter(!showAdvancedFilter)}
                        >
                            <i className="fa fa-angle-down me-1"></i>
                            Show advanced filters
                        </span>
                    }
                    {
                        showAdvancedFilter &&
                        <span
                            className="cursor-pointer text-muted"
                            onClick={() => setShowAdvancedFilter(!showAdvancedFilter)}
                        >
                            <i className="fa fa-angle-up me-1"></i>
                            Hide advanced filters
                        </span>
                    }
                </div>
            </div>
        </form>
    )
}

export default AdvancedFilter