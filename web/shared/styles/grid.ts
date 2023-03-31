export const defineGridStyle = (length: number) => {
    let styles = gridStyle;
    return styles += (length <= 3 ? " grid-cols-3 "  : setResponsiveGrid);
};
export const gridStyle = " grid grid-flow-row mx-auto mt-3 gap-4 grid-cols-4"


const setResponsiveGrid = " max-sm:grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 max-xl:grid-cols-4 "