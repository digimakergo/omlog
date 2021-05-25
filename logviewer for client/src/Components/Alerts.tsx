import styled from "styled-components";

type PropTypes = {
  type: "error" | "info" | "warning" | "debug";
};

export const Alert = ({ type }: PropTypes) => {
  // assign desired colors for alerts based on the alert type
  const color =
    type === "warning"
      ? "var(--yellow)"
      : type === "error"
      ? "var(--red)"
      : type === "info"
      ? "var(--green)"
      : "#333131";

  return <AlertBox color={color}>{type}</AlertBox>;
};

const AlertBox = styled.div<{ color: string; size?: string }>`
  padding: 0.3rem;
  min-width: 70px;
  text-align: center;
  color: #fff;
  font-weight: 500;
  background: ${(props) => props.color};
  font-size: ${(props) => props.size ?? "14px"};
  width: fit-content;
  border-radius: 5px;
`;
