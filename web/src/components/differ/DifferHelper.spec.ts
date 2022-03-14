import {getIndicesOf, searchMatches} from './DifferHelper'
import {Hunk} from "../../__generated__/types";

describe('DifferHelper', () => {
  it('get indexes length', () => {
    const str1 = `diff --git "a/readme.md" "b/readme.md"
index 48533da..f627f7d 100644
--- "a/readme.md"
+++ "b/readme.md"
@@ -1,1 +1,7 @@
-# testing1
\\ no newline at end of file
+# testing1
+dsada
+dasda
+testing2
+das
\\ no newline at end of file
+testing3
+testing4
\\ no newline at end of file
`
    expect(getIndicesOf('testing', str1, false).length).toEqual(5);
    expect(getIndicesOf('da', str1, false).length).toEqual(4);
    expect(getIndicesOf('sa', str1, false).length).toEqual(1);
    expect(getIndicesOf('f', str1, false).length).toEqual(0);

    const n: Hunk = {
        isApplied: false,
        isDismissed: false,
        isOutdated: false,
        id: "d320fd593ace289810cb0991e437fec054ef34d4215bf8e07a429c17fa1fea36",
        patch: str1};
      const m = new Map<string, number[]>();
      m.set(n.id, getIndicesOf('testing', n.patch, false));
      expect(searchMatches(m, [n]).size).toEqual(5);
  })
})
