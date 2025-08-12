import React from 'react';
import Users from './components/Users';
import Todos from './components/Todos';

function App() {
  console.log('Node API base:', process.env.REACT_APP_NODE_API_BASE);
  console.log('Go API base:', process.env.REACT_APP_GO_API_BASE);

  return (
    <div style={{ maxWidth: 800, margin: 'auto', padding: 20 }}>
      <h1>Microservices App: Node & Go Backend APIs</h1>
      <Users />
      <hr />
      <Todos />
    </div>
  );
}

export default App;
