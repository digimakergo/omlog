import { useState } from "react";
import styled from "styled-components";
import { Link, PrimaryText } from "./Styles";
import { BsFunnel, BsCaretDownFill } from "react-icons/bs";

// all available filters & options
type FilterType = {
  name: string;
  options: string[];
};
const allFilters = [
  { name: "Category", options: ["Permission", "Total", "Db", ""] },
  { name: "Level", options: ["Warning", "Info", "Error", "Debug"] },
];

type PropsType = {
  activeFilters: FilterType[];
  setActiveFilters: React.Dispatch<React.SetStateAction<FilterType[]>>;
};

export const Filters = ({ activeFilters, setActiveFilters }: PropsType) => {
  const [expandedFilter, setExpandedFilter] = useState("Level");

  // adds selected filter to active filters
  const addFilter = (filterName: string, filterOption: string) => {
    const filter = activeFilters.find((filter) => filter.name === filterName);
    const modifiedFilterOptions = filter?.options.filter(
      (option) => option !== filterOption
    );
    if (filter) {
      if (filter.options.includes(filterOption)) {
        if (!modifiedFilterOptions) return;
        setActiveFilters((prev) => [
          ...prev.filter((filter) => filter.name !== filterName),
          { name: filterName, options: modifiedFilterOptions },
        ]);
      } else {
        setActiveFilters((prev) => [
          ...prev.filter((f) => f.name !== filterName),
          { name: filterName, options: [...filter.options, filterOption] },
        ]);
      }
    } else {
      setActiveFilters((prev) => [
        ...prev,
        { name: filterName, options: [filterOption] },
      ]);
    }
  };
  return (
    <FilterContainer>
      <Filter>
        <BsFunnel />
        FILTERS
      </Filter>
      {allFilters.map((filter) => {
        const isExpanded = filter.name === expandedFilter;
        return (
          <>
            <Filter
              key={filter.name}
              onClick={() =>
                setExpandedFilter((prev) =>
                  prev === filter.name ? "" : filter.name
                )
              }
            >
              <Icon active={isExpanded}>
                <Link>
                  <BsCaretDownFill />
                </Link>
              </Icon>
              {filter.name}
            </Filter>
            {isExpanded &&
              filter.options.map((option) => {
                const isActive = activeFilters
                  .find((f) => f.name === filter.name)
                  ?.options.includes(option);
                return (
                  <Filter key={option} style={{ marginLeft: "1rem" }}>
                    <input
                      type="checkbox"
                      checked={isActive}
                      onChange={() => addFilter(filter.name, option)}
                    />
                    {option || "Undefined"}
                  </Filter>
                );
              })}
          </>
        );
      })}
    </FilterContainer>
  );
};

const FilterContainer = styled.div`
  max-width: 100%;
  display: flex;
  flex-direction: column;
  padding: 0.5rem;
  border: 1px solid var(--gray);
  background: #fff;
`;

const Filter = styled(PrimaryText)`
  display: grid;
  grid-template-columns: 16px 1fr;
  grid-gap: 0.5rem;
  margin-top: 1rem;
  align-items: center;
  cursor: pointer;
`;

const Icon = styled.div<{ active?: boolean }>`
  transition: transform 0.5s ease;
  transform: ${(props) =>
    props.active ? "rotateZ(360deg)" : "rotateZ(-90deg)"};
`;
