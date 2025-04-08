// Utility functions for the MCP Server

/**
 * Convert a snake_case string to PascalCase
 * @param {string} str - Snake case string
 * @returns {string} Pascal case string
 */
function toPascalCase(str) {
  return str
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join('');
}

/**
 * Check if a SQL query is read-only
 * @param {string} query - SQL query to check
 * @returns {boolean} True if the query is read-only
 */
function isReadOnlyQuery(query) {
  const lowerQuery = query.toLowerCase().trim();
  return lowerQuery.startsWith('select') && 
         !lowerQuery.includes('delete') && 
         !lowerQuery.includes('insert') && 
         !lowerQuery.includes('update') && 
         !lowerQuery.includes('drop') && 
         !lowerQuery.includes('alter') && 
         !lowerQuery.includes('create');
}

module.exports = {
  toPascalCase,
  isReadOnlyQuery
};
