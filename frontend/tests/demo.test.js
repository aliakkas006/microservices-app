const dummyUsers = [
  { id: 1, name: 'Ali Akkas', email: 'ali@example.com' },
  { id: 2, name: 'Rahim', email: 'rahim@example.com' },
  { id: 3, name: 'Karim', email: 'karim@example.com' },
];

describe('Dummy User Data', () => {
  test('should have 3 users', () => {
    expect(dummyUsers.length).toBe(3);
  });
});
