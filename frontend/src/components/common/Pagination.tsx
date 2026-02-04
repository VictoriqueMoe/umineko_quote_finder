import { PAGE_SIZE } from "../../api/endpoints";

interface PaginationProps {
    total: number;
    offset: number;
    onPaginate: (newOffset: number) => void;
}

export function Pagination({ total, offset, onPaginate }: PaginationProps) {
    if (total <= PAGE_SIZE) {
        return null;
    }

    const totalPages = Math.ceil(total / PAGE_SIZE);
    const currentPage = Math.floor(offset / PAGE_SIZE) + 1;
    const hasPrev = offset > 0;
    const hasNext = offset + PAGE_SIZE < total;

    return (
        <div className="pagination">
            <button
                className="pagination-btn"
                disabled={!hasPrev}
                onClick={() => {
                    if (hasPrev) {
                        onPaginate(offset - PAGE_SIZE);
                    }
                }}
            >
                {"\u25C0 Previous"}
            </button>
            <span className="pagination-info">
                Page <span>{currentPage}</span> of <span>{totalPages}</span>
            </span>
            <button
                className="pagination-btn"
                disabled={!hasNext}
                onClick={() => {
                    if (hasNext) {
                        onPaginate(offset + PAGE_SIZE);
                    }
                }}
            >
                {"Next \u25B6"}
            </button>
        </div>
    );
}
