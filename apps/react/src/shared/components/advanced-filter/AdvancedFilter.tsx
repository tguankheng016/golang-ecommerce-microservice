import { InputText } from 'primereact/inputtext';
import { Dispatch, SetStateAction, useState } from 'react'

export interface AdvancedFilterProps {
    filterText: string;
    setFilterText: Dispatch<SetStateAction<string>>;
    children?: React.ReactNode;
}

const AdvancedFilter = ({ filterText, setFilterText, children }: AdvancedFilterProps) => {
    const [showAdvancedFilter, setShowAdvancedFilter] = useState(false);

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
                                name="filterText"
                                placeholder="Search ..."
                                value={filterText}
                                onChange={(e) => setFilterText(e.target.value)}
                            />
                        </div>
                    </div>
                </div>
            </div>
            {
                showAdvancedFilter &&
                <div className="row">
                    {children}
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
                            <i className="fa fa-angle-down"></i>
                            Show advanced filters
                        </span>
                    }
                    {
                        showAdvancedFilter &&
                        <span
                            className="cursor-pointer text-muted"
                            onClick={() => setShowAdvancedFilter(!showAdvancedFilter)}
                        >
                            <i className="fa fa-angle-up"></i>
                            Hide advanced filters
                        </span>
                    }
                </div>
            </div>
        </form>
    )
}

export default AdvancedFilter