import { useState, useEffect } from "react";
import styled from "styled-components";
import {
  Navbar,
  Container,
  FlexBox,
  Link,
  PrimaryText,
  SecondaryText,
  Alert,
  SearchBar,
  Filters,
} from "../../Components";
import Highlight from "react-highlighter";
import ReactPaginate from "react-paginate";
import socketIoClient from "socket.io-client";
import useWebSocket from "react-use-websocket";

//Icons
import { IoIosRocket, IoIosArrowRoundDown } from "react-icons/io";
import { format, differenceInDays } from "date-fns";
import { VscChevronRight } from "react-icons/vsc";
import { CgCloseO } from "react-icons/cg";

// fetch data from backend
import { useData } from "../../Utils";

//WebSocket
const SOCKET_SERVER = "ws://:6001/ws/debug-logs";
//const SOCKET_SERVER = "ws://91.192.221.250:6001/ws/debug-logs";

type LiveData = {
  Category: string;
  DebugId: string;
  Ip: string;
  Level: "error" | "info" | "warning" | "debug";
  Msg: string;
  RequestId: string;
  Time: string;
  Type: string;
  Uri: string;
  UserId: number;
};

export const Dashboard = () => {
  // state holds search keyword to mark the text.
  const [searchKeyword, setSearchKeyword] = useState("");
  const { data, setPage, totalPages } = useData();
  const [liveData, setLiveData] = useState<LiveData[]>([]);
  const [isLiveMode, setLiveMode] = useState(false);

  // state to hold all active filters.
  // by default Level filter is selected, and "error", "info", "warning", "debug" options are checked.

  const [activeFilters, setActiveFilters] = useState([
    { name: "Level", options: ["Warning", "Info", "Error", "Debug"] },
    { name: "Category", options: ["", "Db"] },
  ]);

  // function to remove selected filter from active filters
  const removeSelectedFilter = (filterName: string, filterOption: string) => {
    const filter = activeFilters.find((filter) => filter.name === filterName);
    if (!filter) return;
    const modifiedFilterOptions = filter.options.filter(
      (option) => option !== filterOption
    );
    setActiveFilters((prev) => [
      ...prev.filter((filter) => filter.name !== filterName),
      { name: filterName, options: modifiedFilterOptions },
    ]);
  };

  const { lastJsonMessage } = useWebSocket(SOCKET_SERVER);
  console.log(lastJsonMessage);
  useEffect(() => {
    if (isLiveMode && lastJsonMessage) {
      setLiveData((prev) => [...prev, { ...lastJsonMessage.Logs }]);
    }
  }, [lastJsonMessage, isLiveMode]);

  // format `Date` to time and date string.

  const handleDateTime = (inputDate: string) => {
    const newDate = new Date(inputDate.split(" CEST")[0]);
    const time = format(newDate, "h:mm:ss a");
    const date = format(newDate, "M/dd/yyyy");
    return [time, date];
  };

  return (
    <>
      <Navbar />
      <Layout>
        <Box>
          <PrimaryText size="16px">Digimaker Log Viewer</PrimaryText>
          <Link size="13px" onClick={() => setLiveMode((prev) => !prev)}>
            <IoIosRocket />
            &nbsp; {isLiveMode ? "STOP LIVE MODE" : "START LIVE MODE"}
          </Link>
        </Box>
        {data.length > 1 && (
          <Box>
            {/*
            <SecondaryText size="12px">
              {!isLiveMode && (
                <>
                  {format(
                    new Date(data[0].time.split(" CEST")[0]),
                    "MMM dd, h:mm a"
                  )}{" "}
                </>
              )}
                  </SecondaryText>*/}

            <Link>
              {isLiveMode ? (
                <span style={{ color: "red" }}>*Live*</span>
              ) : (
                <>
                  Last{" "}
                  {differenceInDays(
                    new Date(),
                    new Date(data[0].time.split(" CEST")[0])
                  )}{" "}
                  days
                </>
              )}
            </Link>

            
          </Box>
        )}

        <FilterBlock>
          <Filter>
            <PrimaryText>FILTERS: &nbsp;</PrimaryText>
            {activeFilters.map((filter, index) => {
              if (filter.options.length === 0) return;
              return (
                <Filter key={filter.name + index}>
                  <FilterName>{filter.name}</FilterName>
                  {filter.options.map((option) => (
                    <FilterOption key={option}>
                      {option || "Undefined"}&nbsp;
                      <Link
                        onClick={() =>
                          removeSelectedFilter(filter.name, option)
                        }
                      >
                        <CgCloseO />
                      </Link>
                    </FilterOption>
                  ))}
                </Filter>
              );
            })}
          </Filter>

          <Link width="60px" onClick={() => setActiveFilters([])}>
            Clear All
          </Link>
        </FilterBlock>

        <Section>
          <Filters
            activeFilters={activeFilters}
            setActiveFilters={setActiveFilters}
          />
          <InfoContainer>
            <Box>
              <Sort>
                <Link size="18px">
                  <IoIosArrowRoundDown />
                </Link>
                <Link size="12px">DateTime</Link>
              </Sort>

              <SearchBar setSearchKeyword={setSearchKeyword} />
            </Box>
            {isLiveMode ? (
              <>
                {liveData
                  .filter((d) =>
                    activeFilters
                      .find((f) => f.name === "Level")
                      ?.options.includes(
                        d.Level.slice(0, 1).toUpperCase() + d.Level.slice(1)
                      )
                  )
                  .filter((d) =>
                    activeFilters
                      .find((f) => f.name === "Category")
                      ?.options.includes(
                        d.Category.slice(0, 1).toUpperCase() +
                          d.Category.slice(1)
                      )
                  )
                  ?.reverse()
                  .slice(0, 11)
                  .map((info, index) => {
                    const [time, date] = handleDateTime(info.Time);
                    if (index)
                      return (
                        <InfoTab
                          striped={index % 2 === 0}
                          key={info.DebugId + index}
                        >
                          <InfoTabGroup>
                            <SecondaryText size="13px">{date}</SecondaryText>
                            <PrimaryText size="12px">{time}</PrimaryText>
                          </InfoTabGroup>

                          <SecondaryText size="13px">
                            <Highlight search={searchKeyword}>
                              {info.Msg}
                            </Highlight>
                          </SecondaryText>

                          <AlertBox>
                            <Alert type={info.Level} />
                            <Link size="18px">
                              <VscChevronRight />
                            </Link>
                          </AlertBox>
                        </InfoTab>
                      );
                  })}
              </>
            ) : (
              <>
                {data
                  .filter((d) =>
                    activeFilters
                      .find((f) => f.name === "Level")
                      ?.options.includes(
                        d.level.slice(0, 1).toUpperCase() + d.level.slice(1)
                      )
                  )
                  .filter((d) =>
                    activeFilters
                      .find((f) => f.name === "Category")
                      ?.options.includes(
                        d.category.slice(0, 1).toUpperCase() +
                          d.category.slice(1)
                      )
                  )
                  ?.map((info, index) => {
                    const [time, date] = handleDateTime(info.time);
                    return (
                      <InfoTab striped={index % 2 === 0} key={info.id + index}>
                        <InfoTabGroup>
                          <SecondaryText size="13px">{date}</SecondaryText>
                          <PrimaryText size="12px">{time}</PrimaryText>
                        </InfoTabGroup>

                        <SecondaryText size="13px">
                          <Highlight search={searchKeyword}>
                            {info.msg}
                          </Highlight>
                        </SecondaryText>

                        <AlertBox>
                          <Alert type={info.level} />
                          <Link size="18px">
                            <VscChevronRight />
                          </Link>
                        </AlertBox>
                      </InfoTab>
                    );
                  })}
              </>
            )}

            {!isLiveMode && (
              <Pagination>
                <ReactPaginate
                  pageCount={totalPages}
                  pageRangeDisplayed={2}
                  marginPagesDisplayed={1}
                  onPageChange={(value) => {
                    setPage(value.selected + 1);
                  }}
                  containerClassName={"container"}
                />
              </Pagination>
            )}
          </InfoContainer>
        </Section>
      </Layout>
    </>
  );
};

// styled components

const Pagination = styled.div`
  margin: 1rem;

  display: grid;
  place-items: center;

  .container {
    border: 1px solid var(--gray);
    padding: 1rem;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(30px, 0.15fr));
    width: 100%;
    list-style-type: none;
    grid-gap: 10px;
    color: #fff;
  }
  li a {
    background: var(--primary);
    padding: 5px 10px;
    cursor: pointer;
    border-radius: 5px;
  }

  .previous {
    display: none;
  }

  .selected a {
    background: var(--gray);
    color: #000;
  }
`;

const InfoContainer = styled.div`
  width: 100%;
  padding: 0.5rem;
  display: flex;
  flex-flow: column;
  background: #fff;
  border: 1px solid var(--gray);
`;

const AlertBox = styled.div`
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const InfoTab = styled.div<{ striped?: boolean }>`
  display: grid;
  grid-gap: 0.5rem;
  align-items: center;
  grid-template-columns: 80px 1fr 90px;
  background: ${(props) => (props.striped ? "var(--gray)" : "#fff")};
  padding: 0.5rem;
  width: 100%;
`;

const InfoTabGroup = styled.div`
  display: grid;
  height: 100%;
  grid-template-rows: 1rem 1rem;
  grid-gap: 0.2rem;
`;

const Layout = styled(Container)`
  padding: 0.5rem;
  position: relative;
  margin-top: 40px;
`;

const Section = styled.div`
  display: grid;
  grid-template-columns: 200px 1fr;
  @media (max-width: 700px) {
    grid-template-columns: 1fr;
  }
`;

const Box = styled(FlexBox)`
  padding: 0.5rem;
  border: 1px solid var(--gray);
  background: #fff;
`;

const Sort = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #d1c8c8;
  padding: 0.2rem 0.5rem;
`;

const FilterBlock = styled(Box)`
  @media (max-width: 500px) {
    flex-direction: column;
    align-items: flex-start;
  }
`;

const Filter = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
`;

const FilterName = styled.div`
  background: var(--primary-dark);
  color: #fff;
  font-weight: 500;
  font-size: 13px;
  padding: 3px 5px;
`;

const FilterOption = styled.div`
  color: var(--text-seconday);
  background: var(--primary-light);
  font-size: 13px;
  padding: 3px;
  display: grid;
  grid-template-columns: 1fr 16px;
  align-items: center;
`;
