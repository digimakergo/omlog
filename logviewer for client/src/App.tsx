import "./App.css";
import { createGlobalStyle } from "styled-components";
import { Dashboard } from "./Pages";

const GlobalStyles = createGlobalStyle`
html, body, #root {
  height: 100%;
  
}

:root {
  font-family: 'Montserrat', sans-serif;
  --primary: #195d81;
  --primary-light: #bbdbe7;
  --primary-dark: #152f45;
  --text-primary: #333231;
  --text-secondary: #696868;
  --blue: #0b81e4;
  --gray: #f5f5f5;
  --orange: #fd5203;
  --brown: #7f0001;
  --green: #21a81c;
  --yellow: #fec32e;
  --red: #d30915;
  
}
* {
  box-sizing: border-box;
  font-family: 'Montserrat', sans-serif;
}

body {
  background: var(--gray)
}
`;

function App() {
  return (
    <>
      <GlobalStyles />
      <Dashboard />
    </>
  );
}

export default App;
