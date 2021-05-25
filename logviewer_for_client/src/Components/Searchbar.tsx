import React from "react";
import styled from "styled-components";
import { Link } from "./Styles";
import { CgCloseO } from "react-icons/cg";
import { BiSearchAlt } from "react-icons/bi";
import { useState } from "react";

type PropTypes = {
  width?: string;
  setSearchKeyword: React.Dispatch<React.SetStateAction<string>>;
};

export const SearchBar = ({ width, setSearchKeyword }: PropTypes) => {
  const [inputText, setInputText] = useState("");
  return (
    <Group>
      <Searchbar>
        <Input
          placeholder="Search here..."
          onChange={(e) => setInputText(e.target.value)}
        />
        <Cross>
          <Link
            size="16px"
            onClick={() => {
              setSearchKeyword("");
            }}
          >
            <CgCloseO />
          </Link>
        </Cross>
      </Searchbar>
      <Search>
        <Link size="24px" onClick={() => setSearchKeyword(inputText)}>
          <BiSearchAlt />
        </Link>
      </Search>
    </Group>
  );
};

const Group = styled.div<{ width?: string }>`
  max-width: ${(props) => props.width ?? "300px"};
  display: grid;
  grid-template-columns: 1fr 30px;
`;

const Searchbar = styled.div`
  position: relative;
  width: 100%;
  padding: 0.5rem;
`;

const Cross = styled.div`
  position: absolute;
  top: 0.75rem;
  right: 12px;
`;

const Search = styled.div`
  display: grid;
  place-items: center;
  width: 100%;
`;

const Input = styled.input`
  color: var(--text-primary);
  width: 100%;
  font-size: 14px;
  border: 1px solid #b7afaf;
  padding: 0.2rem 0.5rem;
`;
