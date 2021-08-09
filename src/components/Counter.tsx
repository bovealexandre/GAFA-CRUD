import React, { useState } from 'react';

const Counter = ({ children }) => {
    const [count, setCount] = useState(0);
    const add = () => setCount((i: number) => i + 1);
    const subtract = () => setCount((i: number) => i - 1);

    return (
        <>
            <div className="counter">
                <button onClick={subtract}>-</button>
                <pre>{count}</pre>
                <button onClick={add}>+</button>
            </div>
            <div className="children">{children}</div>
        </>
    );
}

export default Counter