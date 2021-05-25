import React from "react";
import styled from "styled-components";
import { Container, FlexBox } from "./Styles";
import mainlogo from "./Icon2.png"; 

export const Navbar = () => {
  return (
    <Nav>
      <Container>
        <FlexBox>
          <Logo>
            <img src={mainlogo}width="23" height="23  "/> &nbsp;Digimaker
          </Logo>
          <NavItem>Dashboard</NavItem>
        </FlexBox>
      </Container>
    </Nav>
  );
};


const Nav = styled.div`
  width: 100%;
  background: var(--text-primary);
  display: flex;
  padding: 0.5rem 1rem;
  position: fixed;
  top: 0;
  z-index: 100;
`;

const Logo = styled.div`
  color: #fff;
  font-weight: 300;
  
`;

const NavItem = styled.div`
  color: #fff;
  font-weight: 500;
`;
