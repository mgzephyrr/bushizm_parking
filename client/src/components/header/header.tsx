import React from "react";

interface IHeaderProps {
    children?: React.ReactNode;
}

export default function Header({children}: IHeaderProps) {
    return (
        <div style={{
            boxSizing: "border-box",
            backgroundColor: "#569EB7",
            height: "100px",
        }}>
            {children}
        </div>
    )
}