---
name: Task
about: Create a new task for the Innovation team to work on
title: ''
assignees: ''
---

<details>
  <summary>
    $\textcolor{Darkorange}{\textsf{DoR}}$ 
  </summary>
  Tasks can’t be started if the following info doesn’t exist.

* Add an estimation to the ticket using dedicated field
* If needed, add a time-box validated by PM for all Investigation / exploration tickets
* Think and write down what you'll do to achieve this ticket. Add an action plan in comment
  * An action plan is:
  * An ordered list of tasks
        * with an owner on each task
        * and an estimation
        * if blockers or dependencies are identified, they must be clearly listed and a solution should be found before starting the issue.
* For front-end tasks, designs must be linked, accessible by all and MR must contain [context], [before/after image][reproduction instructions]

If one of this point or info is missing, please raise the point ASAP to the PM.

</details>



**Context**

*Describe / explain why we should do this: motivations, context or other info. Keep it brief and simple*


**User flow**

*Describe the user flow using user stories so the end result is super clear*


**How to**

*List the step-by-step to get it do if needed*


**Technical details**

*Give the technical insights so anyone in the team can tackle the tasks*

**QA testing**

*Does this task require some QA tests ?*
*If yes, explain how to validate it*

<details>
  <summary>$\textcolor{Red}{\textsf{DoD}}$ </summary>
Before putting this ticket in code review, tick all the boxes bellow.
  More details [here] (https://www.notion.so/massa-innoteam/Plan-for-the-mainnet-c574da44a4854eb3841a5f2e93a2977c?pvs=4#e7db6fa53fa84264954075011432ce70) & [here] (https://www.notion.so/massa-innoteam/Industrialization-of-Frontend-0f7425f14cd3490a949f31978916ee41?pvs=4) if needed
  
- [ ] Pull request is small and approved by 2 reviewers
    - Max of 10 files updated and 500 lines of code added
    - Break things down as much as you can
- [ ] You are proud of what will / has been merged
- [ ] Code and functionality implemented is working on all OS
    - Windows 10 +
    - Linux Ubuntu
    - MacOS catalina +
- [ ] Endpoints are covered by units tests and are monitored (ie: we are notified -somehow- when something is down)
- [ ] Implemented screens are pixel perfect with the designs for the following screen sizes
    - 1920 x 1080 px
    - 1440×900 px
    - 1366 x 768 px
- [ ] All info must appear in less than 1 sec on the front-end (when applicable)
- [ ] Functionalities are fully working (errors messages exist, all use cases are covered - when applicable)
- [ ] Related documentation has been updated if needed
- [ ] Functionality are QA reviewed on every OS
</details>
