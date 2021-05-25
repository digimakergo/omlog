import React from "react";
import styled from "styled-components";

export const Container = styled.div`
  max-width: 1440px;
  width: 100%;
  margin: auto;
`;
export const FlexBox = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

export const Link = styled.a<{ size?: string; width?: string }>`
  text-decoration: none;
  color: var(--blue);
  font-weight: 500;
  font-size: ${(props) => props.size ?? "14px"};
  display: flex;
  align-items: center;
  cursor: pointer;
  width: ${(props) => props.width ?? "fit-content"};
`;

export const PrimaryText = styled.div<{ size?: string }>`
  color: var(--text-primary);
  font-weight: 600;
  font-size: ${(props) => props.size ?? "14px"};
`;
export const SecondaryText = styled.div<{ size?: string }>`
  color: var(--text-secondary);
  font-weight: 400;
  font-size: ${(props) => props.size ?? "14px"};
  word-break: break-word;
`;
