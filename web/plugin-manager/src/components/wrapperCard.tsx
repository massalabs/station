import React from "react";

type Props = {
    children: any;
};

const wrapperCard = (props: Props) => {
    return <div>{props.children}</div>;
};

export default wrapperCard;
